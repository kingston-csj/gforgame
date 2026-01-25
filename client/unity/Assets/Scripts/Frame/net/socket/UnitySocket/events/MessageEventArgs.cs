using System;
using Nova.Net.Socket;


namespace Nova.Net.UnitySocket
{
    public class MessageEventArgs : EventArgs
    {
        public SocketDataFrame DataFrame;

        internal MessageEventArgs(SocketDataFrame dataFrame)
        {
            DataFrame = dataFrame;
        }
    }
}