using System;
using System.Collections;
using System.Collections.Generic;
using Frame.Codec.Bean;
using Frame.Commons.Utils;
using Nova.Commons.Util;
using Nova.Net.Socket;

namespace Nova.Codec.bean
{
    /// <summary>
    /// key 强制为 string 类型
    /// 元素Value可以是父类或抽象类
    /// </summary>
    public class MapCodec2 : Codec
    {
        public override object Decode(ByteBuff buff, Type type)
        {
            // 1. 校验目标类型：必须是 Dictionary<string, TValue> 类型
            if (!type.IsGenericType || type.GetGenericTypeDefinition() != typeof(Dictionary<,>))
            {
                throw new ArgumentException($"目标类型必须是 Dictionary<TKey, TValue>，实际类型：{type.FullName}");
            }

            Type[] genericArgs = type.GetGenericArguments();
            Type keyType = genericArgs[0];
            Type valueType = genericArgs[1];
            // 强制 Key 为 string（符合业务约定）
            if (keyType != typeof(string))
            {
                throw new ArgumentException($"Map Key 类型必须是 string，实际类型：{keyType.FullName}");
            }

            // 读取 Map 长度（short 类型）
            short size = buff.ReadShort();
            if (size < 0)
            {
                throw new InvalidOperationException($"Map 长度不能为负数：{size}");
            }

            // 处理空字典：直接返回空实例
            if (size == 0)
            {
                return Activator.CreateInstance(type);
            }

            // 非空字典：初始化目标字典 
            IDictionary targetDict = (IDictionary)Activator.CreateInstance(type);
            // 读取状态位
            byte status = buff.ReadByte();

            // 循环解码键值对
            for (int i = 0; i < size; i++)
            {
                // 解码 Key（固定为 string）
                string key = (string)GetCodec(typeof(string)).Decode(buff, typeof(string));
                // 确定 Value 的实际解码类型
                Type eleType = valueType; // 默认用泛型指定的 Value 类型（如 BaseVo）
                if (status == 1)
                {
                    int messageId = buff.ReadInt();
                    eleType = BeanCodecContext.MessageFactory.GetMessageType(messageId);
                    // 校验动态类型是否兼容目标 Value 类型（如 BaseVo 子类）
                    if (!valueType.IsAssignableFrom(eleType))
                    {
                        throw new InvalidCastException(
                            $"动态解析的类型 {eleType.FullName} 无法赋值给目标类型 {valueType.FullName}");
                    }
                }

                // 解码 Value（使用实际类型编解码器）
                Codec valueCodec = GetCodec(eleType);
                object value = valueCodec.Decode(buff, eleType);
                // 将键值对添加到目标字典
                targetDict.Add(key, value);
            }

            // 返回目标类型的字典（如 Dictionary<string, BaseVo>）
            return targetDict;
        }

        public override void Encode(ByteBuff buff, object value)
        {
            if (value == null)
            {
                // 空 Map 写入长度 0
                buff.WriteShort(0);
                return;
            }

            var nonGenericDict = value as IDictionary;
            var map = new Dictionary<string, object>();

            ISet<Type> set = new HashSet<Type>();
            foreach (DictionaryEntry entry in nonGenericDict)
            {
                // 强制校验 Key 是 string（符合你的业务约定）
                if (!(entry.Key is string key))
                {
                    throw new ArgumentException($"字典 Key 必须是 string 类型，实际 Key 类型：{entry.Key.GetType().FullName}");
                }

                // Value 直接转 object（BaseVo 实例会被保留，后续编码可识别）
                map.Add(key, entry.Value);
                set.Add(entry.Value.GetType());
            }

            // 写入 Map 长度
            short size = (short)map.Count;
            buff.WriteShort(size);
            if (size == 0)
            {
                return;
            }

//        key统一为string，根据value类型判断状态
//        1:基本类型，写入状态0;
//        2:不是基本类型，且元素类型是一样的，且不是抽象类(接口)写入状态0
//        3:否则，写入状态1，然后在迭代的时候，同时写入每个元素的类型id
            byte status = 0;
            Type dictType = value.GetType();
            Type[] genericArgs = dictType.GetGenericArguments();
            // 0是key类型，1是value类型
            if (!TypeUtil.IsPrimitiveOrString(genericArgs[1]))
            {
                if (set.Count > 1 || genericArgs[1].IsAbstract || genericArgs[1].IsInterface)
                {
                    status = 1;
                }
            }

            buff.WriteByte(status);
            // 循环编码键值对
            foreach (var entry in map)
            {
                string key = entry.Key;
                object valueObj = entry.Value;

                // 编码 Key（强制用 StringCodec）
                GetCodec(typeof(string)).Encode(buff, key);

                if (status == 1)
                {
                    // 状态为1，写入value的类型id
                    int messageId = BeanCodecContext.MessageFactory.GetMessageCmd(valueObj.GetType());
                    buff.WriteInt(messageId);
                }

                Codec valueCodec = GetCodec(valueObj.GetType());
                valueCodec.Encode(buff, valueObj);
            }
        }
    }
}