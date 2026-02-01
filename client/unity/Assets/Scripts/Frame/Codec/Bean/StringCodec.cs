using System;
using System.Text;
using Nova.Commons.Util;

namespace Nova.Codec.bean
{
    /// <summary>
    /// 字符串编解码器
    /// </summary>
    public class StringCodec : Codec
    {
        public override object Decode(ByteBuff buff, Type type)
        {
            int length = buff.ReadInt();
            if (length == 0)
                return null;
            byte[] bytes = buff.ReadBytes(length);
            return Encoding.UTF8.GetString(bytes);
        }

        public override void Encode(ByteBuff buff, object value)
        {
            if (value == null)
            {
                buff.WriteInt(0);
                return;
            }

            string str = value.ToString();
            byte[] bytes = Encoding.UTF8.GetBytes(str);
            buff.WriteInt(bytes.Length);
            buff.WriteBytes(bytes, bytes.Length);
        }
    }
    
}