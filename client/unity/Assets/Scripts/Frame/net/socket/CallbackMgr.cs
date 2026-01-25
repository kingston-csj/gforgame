using System.Collections.Generic;

namespace Nova.Net.Socket
{
    public class CallbackMgr
    {
        /// <summary>
        /// todo 定时检测，并删除超时的请求
        /// </summary>
        private static Dictionary<int, MessageCallback> callbacks = new();

        public static void Register(int id, MessageCallback callback)
        {
            callbacks[id] = callback;
        }

        /// <summary>
        ///     获取回调并删除
        /// </summary>
        /// <param name="id">回调ID</param>
        /// <returns>回调</returns>
        public static MessageCallback Fetch(int id)
        {
            if (callbacks.TryGetValue(id, out MessageCallback callback))
            {
                // 删除
                callbacks.Remove(id);
                return callback;
            }

            return null;
        }
    }
    
}