using System;
using UnityEngine.Serialization;

namespace Nova.Net.Socket
{
    [Serializable]
    public class SocketDataFrame
    {
        /// <summary>
        /// 消息包序列(监听方法指针或请求方法指针，0为服务器主动推送)
        /// </summary>
        public int index;

        public int size;

        /// <summary>
        /// 消息 cmd
        /// </summary>
        public int cmd;

        /// <summary>
        /// 消息json格式
        /// </summary>
        public string msgJson;


        /// <summary>
        /// 消息真正的对象
        /// </summary>
        public Message message;
        
        /// <summary>
        /// 消息原生数据
        /// </summary>
        public byte[] rawData;

    }
}