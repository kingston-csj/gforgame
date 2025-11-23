using System;
using Game.Net.Message.Hero;
using Nova.Net.Socket;

namespace Game.Net
{
    /// <summary>
    ///     游戏Socket接收管理类
    /// </summary>
    public class MessageRouter
    {
        /// <summary>
        ///     注册消息回调函数， 将消息id与具体的消息回调方法绑定在一起
        /// </summary>
        /// <typeparam name="T">消息数据类型</typeparam>
        /// <param name="resCallBack">消息回调函数，即被MessageHandler标记的方法</param>
        /// <param name="res_cmd">命令ID</param>
        public void RegisterCallbackDelegate<T>(int res_cmd, Action<T> resCallBack) where T : Nova.Net.Socket.Message
        {
            var tem = new MessageCallback(typeof(T), message =>
            {
                resCallBack((T)message);
            });
            MessageDispatcher.Register(res_cmd, tem);
        }
        
        /// <summary>
        ///     所有英雄信息响应
        /// </summary>
        /// <param name="data"></param>
        [MessageHandler]
        public void OnReceivePushAllHeroInfo(PushAllHeroInfo data)
        {
            
        }
        
    }
}