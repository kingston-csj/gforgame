namespace Nova.Net.Socket
{
    public interface SocketIoDispatcher
    {
        /// <summary>
        ///     处理连接成功事件
        /// </summary>
        void OnOpen();

        /// <summary>
        ///     处理连接关闭事件
        /// </summary>
        void OnClose();
        
        /// <summary>
        ///     处理消息事件
        /// </summary>
        /// <param name="dataFrame"> 包含消息数据的SocketDataFrame对象 </param>
        void OnMessage(SocketDataFrame dataFrame);
    }
}