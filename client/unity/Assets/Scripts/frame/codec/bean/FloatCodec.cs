using System;
using Nova.Commons.Util;

namespace Nova.Codec.bean
{
    public class FloatCodec : Codec
    {
        public override object Decode(ByteBuff buff, Type type)
        {
            return buff.ReadFloat();
        }

        public override void Encode(ByteBuff buff, object value)
        {
            if (value == null)
            {
                buff.WriteFloat(0);
                return;
            }

            buff.WriteFloat(Convert.ToSingle(value));
        }
    }
}