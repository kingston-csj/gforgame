using System;
using Nova.Commons.Util;

namespace Nova.Codec.bean
{
    /// <summary>
    /// 数组编解码
    /// 用于编码和解码数组类型的消息
    /// </summary>
    public class ArrayCodec : Codec
    {
        public override object Decode(ByteBuff buff, Type type)
        {
            // 读取数组长度
            int length = buff.ReadShort();
            if (length <= 0)
            {
                return Array.CreateInstance(type.GetElementType(), 0);
            }

            // 创建数组并逐个解码元素
            Array array = Array.CreateInstance(type.GetElementType(), length);
            Codec elementCodec = GetCodec(type.GetElementType());

            for (int i = 0; i < length; i++)
            {
                object element = elementCodec.Decode(buff, type.GetElementType());
                array.SetValue(element, i);
            }

            return array;
        }

        public override void Encode(ByteBuff buff, object value)
        {
            if (value == null)
            {
                buff.WriteShort(0);
                return;
            }

            Array array = value as Array ?? throw new ArgumentException("value 必须是数组类型");
            int length = array.Length;
            buff.WriteShort((short)length);

            if (length == 0)
            {
                return;
            }

            // 逐个编码元素
            Codec elementCodec = GetCodec(array.GetType().GetElementType());
            for (int i = 0; i < length; i++)
            {
                object element = array.GetValue(i);
                elementCodec.Encode(buff, element);
            }
        }
    }
    
}