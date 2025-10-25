using System;
using Nova.Commons.Util;

namespace Nova.Codec.bean
{
    /// <summary>
    /// byte编解码
    /// </summary>
    public class ByteCodec : Codec
    {
        public override object Decode(ByteBuff buff, Type type)
        {
            return buff.ReadByte();
        }

        public override void Encode(ByteBuff buff, object value)
        {
            if (value == null)
            {
                buff.WriteByte(0);
                return;
            }

            buff.WriteByte(Convert.ToByte(value));
        }
    }
    
}