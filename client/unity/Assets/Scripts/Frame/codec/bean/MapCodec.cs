using System;
using System.Collections.Generic;
using Nova.Commons.Util;

namespace Nova.Codec.bean
{
    /// <summary>
    /// key 强制为 string 类型，后续再做升级
    /// </summary>
    public class MapCodec : Codec
    {
        public override object Decode(ByteBuff buff, Type type)
        {
            // 读取 Map 长度（short 类型）
            short size = buff.ReadShort();
            if (size < 0)
            {
                throw new InvalidOperationException($"Map 长度不能为负数：{size}");
            }

            // 创建 Map 实例（处理接口/抽象类）
            Dictionary<string, object> map = new Dictionary<string, object>(size);
            if (size == 0)
            {
                return map;
            }

            Type valueType = typeof(object);
            Codec valueCodec = Codec.GetCodec(valueType);

            // 循环解码键值对（Key 固定为 string）
            for (int i = 0; i < size; i++)
            {
                // 解码 Key（强制用 StringCodec）
                string key = (string)Codec.GetCodec(typeof(string)).Decode(buff, typeof(string));

                // 解码 Value（使用指定的 Value 类型编解码器）
                object value = valueCodec.Decode(buff, valueType);

                map.Add(key, value);
            }

            return map;
        }

        public override void Encode(ByteBuff buff, object value)
        {
            if (value == null)
            {
                // 空 Map 写入长度 0
                buff.WriteShort((short)0);
                return;
            }

            // 强制转换为 IDictionary<string, object>（非此类型会抛异常）
            if (!(value is IDictionary<string, object> map))
            {
                throw new ArgumentException($"Map 类型必须是 IDictionary<string, object>，实际类型：{value.GetType().FullName}");
            }

            // 写入 Map 长度
            short size = (short)map.Count;
            buff.WriteShort(size);
            if (size == 0)
            {
                return;
            }

            // 循环编码键值对
            foreach (var entry in map)
            {
                string key = entry.Key;
                object valueObj = entry.Value;

                // 编码 Key（强制用 StringCodec）
                Codec.GetCodec(typeof(string)).Encode(buff, key);

                // 编码 Value（根据实际类型获取编解码器）
                Type valueType = valueObj?.GetType() ?? typeof(object);
                Codec valueCodec = Codec.GetCodec(valueType);
                valueCodec.Encode(buff, valueObj);
            }
        }
    }
}