using System;

namespace Nova.Net.Socket
{
    /// <summary>
    /// 对于每一个以Req开头的消息，都有一个以Res开头的响应消息。
    /// 此类为响应消息的基类
    /// </summary>
    [Serializable]
    public class Response : Message
    {
        /// <summary>
        /// 错误码 0表示成功，非0表示失败
        /// </summary>
        public int code;
    }
}