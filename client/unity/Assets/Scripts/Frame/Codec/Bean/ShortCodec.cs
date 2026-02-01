using System;
using Nova.Commons.Util;

namespace Nova.Codec.bean
{
    /// <summary>
    /// short编解码
    /// </summary>
    public class ShortCodec : Codec
    {
        public override object Decode(ByteBuff buff, Type type)
        {
            return buff.ReadShort();
        }

        public override void Encode(ByteBuff buff, object value)
        {
            if (value == null)
            {
                buff.WriteShort(0);
                return;
            }

            buff.WriteShort(Convert.ToInt16(value));
        }
    }
    
}