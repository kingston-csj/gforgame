using System;
using Nova.Commons.Util;

namespace Nova.Codec.bean
{
    /// <summary>
    /// bool编解码器
    /// </summary>
    public class BoolCodec : Codec
    {
        public override object Decode(ByteBuff buff, Type type)
        {
            return buff.ReadBool();
        }

        public override void Encode(ByteBuff buff, object value)
        {
            if (value == null)
            {
                buff.WriteBool(false);
                return;
            }

            buff.WriteBool(Convert.ToBoolean(value));
        }
    }
}