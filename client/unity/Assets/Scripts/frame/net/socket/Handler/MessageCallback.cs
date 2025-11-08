using System;

namespace Nova.Net.Socket
{
    public class MessageCallback
    {
        /// <summary>
        /// 消息数据类型
        /// </summary>
        public Type type;

        /// <summary>
        /// 接收数据回调
        /// </summary>
        public Action<Message> callback;

        public MessageCallback(Type type, Action<Message> reqCallBack)
        {
            this.type = type;
            this.callback = reqCallBack;
        }
        
    }
}