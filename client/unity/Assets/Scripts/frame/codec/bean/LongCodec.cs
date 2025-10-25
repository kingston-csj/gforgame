using System;
using Nova.Commons.Util;

namespace Nova.Codec.bean
{
    /// <summary>
    /// long 编解码
    /// </summary>
    public class LongCodec : Codec
    {
        public override object Decode(ByteBuff buff, Type type)
        {
            return buff.ReadLong();
        }

        public override void Encode(ByteBuff buff, object value)
        {
            if (value == null)
            {
                buff.WriteLong(0);
                return;
            }

            buff.WriteLong(Convert.ToInt64(value));
        }
    }
}