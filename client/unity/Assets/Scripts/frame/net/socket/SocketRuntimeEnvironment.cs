using Nova.Codec;
using System;

namespace Nova.Net.Socket
{
    /// <summary>
    /// socket运行时环境
    /// </summary>
    public class SocketRuntimeEnvironment
    {
        private readonly Type _messageRouterType;

        private readonly MessageCodec _messageCodec;

        private readonly MessageFactory _messageFactory;

        /// <summary>
        /// 是否使用二进制帧
        /// </summary>
        private bool _usedBinaryFrame = false;

        public SocketRuntimeEnvironment(Type messageRouterType, MessageCodec messageCodec,
            MessageFactory messageFactory)
        {
            this._messageRouterType = messageRouterType;
            this._messageCodec = messageCodec;
            this._messageFactory = messageFactory;
        }

        /// 自动生成getter和setter
        /// <summary>
        ///     消息路由器类型
        /// </summary>
        public Type MessageRouterType => _messageRouterType;

        /// <summary>
        ///     消息编码器
        /// </summary>
        public MessageCodec MessageCodec => _messageCodec;

        /// <summary>
        ///     消息工厂
        /// </summary>
        public MessageFactory MessageFactory => _messageFactory;

        /// <summary>
        /// 是否使用二进制帧
        /// </summary>
        public bool UsedBinaryFrame
        {
            get => _usedBinaryFrame;
            set => _usedBinaryFrame = value;
        }
        
    }
}