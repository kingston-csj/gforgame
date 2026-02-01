using System;
using System.Collections.Generic;
using System.Reflection;
using Nova.Commons.Util;

namespace Nova.Codec.bean
{
    /// <summary>
    /// 消息体编码解码基类，同时管理所有类型的字段编解码器
    /// 除了基础类型的编码解码，还支持集合类型的编码解码，暂不支持Map类型
    /// 本工具对于整数的编解码，并没有使用不定长编码，而是使用固定长度编码
    /// 如果业务有特殊需要，可以自行实现整数的不定长编解码方式，并通过 Replace 方法替换
    /// </summary>
    public abstract class Codec
    {
        /// <summary>
        /// 类型→编解码器映射（线程安全）
        /// </summary>
        private static readonly Dictionary<Type, Codec> _typeToCodecs =
            new();

        static Codec()
        {
            // 注册基础类型编解码器（与 Java 保持一致）
            Register(typeof(bool), new BoolCodec());
            Register(typeof(bool?), new BoolCodec());
            Register(typeof(byte), new ByteCodec());
            Register(typeof(byte?), new ByteCodec());
            Register(typeof(short), new ShortCodec());
            Register(typeof(short?), new ShortCodec());
            Register(typeof(int), new IntCodec());
            Register(typeof(int?), new IntCodec());
            Register(typeof(float), new FloatCodec());
            Register(typeof(float?), new FloatCodec());
            Register(typeof(double), new DoubleCodec());
            Register(typeof(double?), new DoubleCodec());
            Register(typeof(long), new LongCodec());
            Register(typeof(long?), new LongCodec());
            Register(typeof(string), new StringCodec());
            Register(typeof(List<>), new CollectionCodec());
            Register(typeof(HashSet<>), new CollectionCodec());
            Register(typeof(object[]), new ArrayCodec());
            Register(typeof(Dictionary<,>), new MapCodec());
        }

        /// <summary>
        /// 注册类型与编解码器的绑定（重复注册会抛异常）
        /// </summary>
        public static void Register(Type type, Codec codec)
        {
            if (!_typeToCodecs.TryAdd(type, codec))
            {
                throw new InvalidOperationException($"{type.FullName} 已存在对应的编解码器，请勿重复注册");
            }
        }

        /// <summary>
        /// 替换类型对应的编解码器（支持自定义扩展）
        /// </summary>
        public static void Replace(Type type, Codec codec)
        {
            _typeToCodecs[type] = codec;
        }

        /// <summary>
        /// 获取指定类型的编解码器（自定义Bean会自动生成BeanCodec）
        /// </summary>
        public static Codec GetCodec(Type type)
        {
            // 1. 已注册直接返回
            if (_typeToCodecs.TryGetValue(type, out Codec existingCodec))
            {
                return existingCodec;
            }

            // 2. 数组类型返回数组编解码器
            if (type.IsArray)
            {
                return _typeToCodecs[typeof(object[])];
            }

            // 3. 集合类型（泛型）返回集合编解码器
            if (type.IsGenericType)
            {
                Type genericTypeDef = type.GetGenericTypeDefinition();
                if (genericTypeDef == typeof(List<>) || genericTypeDef == typeof(HashSet<>))
                {
                    return _typeToCodecs[typeof(List<>)];
                }

                // Map 类型（Dictionary/ConcurrentDictionary 等实现 IDictionary<TKey, TValue> 的泛型类）
                if (typeof(IDictionary<,>).IsAssignableFrom(genericTypeDef))
                {
                    // 匹配注册的 Map 编解码器（假设注册时用的键是 typeof(IDictionary<string, object>) 或 typeof(Dictionary<,>)）
                    // 若注册时用的是具体类型（如 Dictionary<string, object>），需对应修改这里的匹配键
                    return _typeToCodecs[typeof(IDictionary<string, object>)];
                }
            }

            // 4. 自定义Bean：解析字段并生成BeanCodec
            List<FieldCodecMeta> fieldMetas = new();
            Type currentType = type;

            // 遍历所有父类（直到Object）
            while (currentType != typeof(object))
            {
                FieldInfo[] fields =
                    currentType.GetFields(BindingFlags.Instance | BindingFlags.NonPublic | BindingFlags.Public |
                                          BindingFlags.DeclaredOnly);
                foreach (FieldInfo field in fields)
                {
                    // 跳过 static/final/transient 字段
                    if (field.IsStatic || field.IsInitOnly ||
                        field.IsDefined(typeof(NonSerializedAttribute), inherit: false))
                    {
                        continue;
                    }

                    // 忽略本地字段
                    if (field.GetCustomAttribute<FieldIgnore>() != null)
                    {
                        continue;
                    }

                    Type fieldType = field.FieldType;
                    Codec fieldCodec = GetCodec(fieldType); // 递归获取字段类型的编解码器

                    fieldMetas.Add(new FieldCodecMeta(field, fieldCodec));
                }

                currentType = currentType.BaseType;
            }

            // 生成Bean编解码器并注册
            Codec beanCodec = BeanCodec.Create(fieldMetas);
            Register(type, beanCodec);
            return beanCodec;
        }

        /// <summary>
        /// 解码：从ByteBuff读取数据并转为指定类型
        /// </summary>
        /// <param name="buff">数据缓冲区</param>
        /// <param name="type">目标类型</param>
        public abstract object Decode(ByteBuff buff, Type type);

        /// <summary>
        /// 编码：将对象写入ByteBuff
        /// </summary>
        /// <param name="buff">数据缓冲区</param>
        /// <param name="value">要编码的对象</param>
        public abstract void Encode(ByteBuff buff, object value);
    }
}