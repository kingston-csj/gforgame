using Nova.Codec;
using System;
using Frame.Commons.Utils;

namespace Nova.Net.Socket
{
    /// <summary>
    /// socket运行时环境
    /// </summary>
    public class SocketRuntimeEnvironment
    {
        private readonly Type _messageRouterType;

        private readonly MessageCodec _messageCodec;

        private readonly IMessageFactory _messageFactory;

        /// <summary>
        /// 是否使用二进制帧
        /// 仅在使用WebSocket时生效
        /// 因为WebSocket同时支持二进制帧以及文本帧
        /// 而Socket只能使用二进制帧
        /// </summary>
        private bool _usedWebSocketBinaryFrame = false;

        public SocketRuntimeEnvironment(Type messageRouterType, MessageCodec messageCodec,
            IMessageFactory messageFactory)
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
        public IMessageFactory MessageFactory => _messageFactory;

        /// <summary>
        /// 是否使用二进制帧
        /// </summary>
        public bool UsedBinaryFrame
        {
            get => _usedWebSocketBinaryFrame;
            set => _usedWebSocketBinaryFrame = value;
        }
    }
}