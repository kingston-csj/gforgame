using System;

namespace Nova.Codec
{
    public interface MessageCodec
    {
        /// <summary>
        /// 解码数据
        /// </summary>
        /// <param name="type"></param>
        /// <param name="data"></param>
        /// <returns></returns>
        Object Decode(Type type, byte[] data);

        /// <summary>
        /// 编码数据
        /// </summary>
        /// <param name="data"></param>
        /// <returns></returns>
        byte[] Encode(Object data);
    }
    
}