namespace Nova.Net.UnitySocket
{
    /// <summary>
    /// Socket连接各种状态
    /// </summary>
    public enum SocketState
    {
        /// <summary>
        /// 初始状态，未建立连接
        /// </summary>
        None = 0,

        /// <summary>
        /// 正在建立连接
        /// </summary>
        Connecting = 1,

        /// <summary>
        /// 连接已建立
        /// </summary>
        Opened = 2,

        /// <summary>
        /// 连接已完全关闭
        /// </summary>
        Closed = 3,
    }
}