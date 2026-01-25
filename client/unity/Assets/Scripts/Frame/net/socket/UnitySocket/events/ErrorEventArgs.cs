using System;

namespace Nova.Net.UnitySocket
{
    public class ErrorEventArgs : EventArgs
    {

        public Exception message;

        internal ErrorEventArgs(Exception message)
        {
            this.message = message;
        }
    }
}
