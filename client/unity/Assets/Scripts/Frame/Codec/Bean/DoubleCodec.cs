using System;
using Nova.Commons.Util;

namespace Nova.Codec.bean
{
    /// <summary>
    /// double 类型的编解码实现
    /// </summary>
    public class DoubleCodec : Codec
    {
        public override object Decode(ByteBuff buff, Type type)
        {
            return buff.ReadDouble();
        }

        public override void Encode(ByteBuff buff, object value)
        {
            if (value == null)
            {
                buff.WriteDouble(0);
                return;
            }

            buff.WriteDouble(Convert.ToDouble(value));
        }
    }
}