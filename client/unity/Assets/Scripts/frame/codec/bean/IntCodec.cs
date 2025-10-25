using System;
using Nova.Commons.Util;

namespace Nova.Codec.bean
{
    /// <summary>
    /// 整数类型编解码器（固定4字节，大端序）
    /// </summary>
    public class IntCodec : Codec
    {
        public override object Decode(ByteBuff buff, Type type)
        {
            return buff.ReadInt();
        }

        public override void Encode(ByteBuff buff, object value)
        {
            if (value == null)
            {
                buff.WriteInt(0);  
                return;
            }

            buff.WriteInt(Convert.ToInt32(value));
        }
    }
}