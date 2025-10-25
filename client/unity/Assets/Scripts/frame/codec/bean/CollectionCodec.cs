using System;
using System.Collections;
using System.Reflection;
using Nova.Commons.Util;

namespace Nova.Codec.bean
{
    /// <summary>
    /// 集合编解码器（支持List/HashSet，格式：[长度(4字节)][元素1][元素2]...）
    /// </summary>
    public class CollectionCodec : Codec
    {
        public override object Decode(ByteBuff buff, Type type)
        {
            // 读取集合长度
            int length = buff.ReadShort();
            if (length <= 0)
            {
                return Activator.CreateInstance(type);
            }

            // 获取集合元素类型（如果wrapper未指定，从泛型参数解析）
            Type elementType = type.GetGenericArguments()[0];
            Codec elementCodec = GetCodec(elementType);

            // 创建集合并逐个添加元素
            ICollection collection = (ICollection)Activator.CreateInstance(type);
            MethodInfo addMethod = type.GetMethod("Add", new[] { elementType });

            for (int i = 0; i < length; i++)
            {
                object element = elementCodec.Decode(buff, elementType);
                addMethod.Invoke(collection, new[] { element });
            }

            return collection;
        }

        public override void Encode(ByteBuff buff, object value)
        {
            if (value == null)
            {
                buff.WriteShort(0);
                return;
            }

            // 关键修复：使用非泛型ICollection接口
            ICollection collection = value as ICollection
                                     ?? throw new ArgumentException(
                                         $"value 必须是集合类型（实现 ICollection 接口），当前类型：{value.GetType().FullName}");

            buff.WriteShort((short)collection.Count);
            if (collection.Count == 0)
            {
                return;
            }

            // 获取集合元素类型 
            Type elementType = GetCollectionElementType(value.GetType());
            Codec elementCodec = GetCodec(elementType);

            // 遍历集合元素编码（非泛型ICollection需强制转换为IEnumerable）
            foreach (object element in (IEnumerable)collection)
            {
                elementCodec.Encode(buff, element);
            }
        }

        /// <summary>
        /// 解析集合的元素类型（泛型集合）
        /// </summary>
        private Type GetCollectionElementType(Type collectionType)
        {
            if (!collectionType.IsGenericType)
            {
                throw new NotSupportedException("仅支持泛型集合类型");
            }

            Type[] genericArgs = collectionType.GetGenericArguments();
            if (genericArgs.Length != 1)
            {
                throw new NotSupportedException("仅支持单泛型参数的集合类型");
            }

            return genericArgs[0];
        }
    }
}