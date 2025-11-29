using UnityEngine;

namespace Nova.Logger
{
    
    public class LoggerUtil
    {
        public static bool Enabled = false;

        /// <summary>
        /// 打印info日志
        /// </summary>
        /// <param name="message"></param>
        public static void Info(string message)
        {
            if (Enabled)
            {
                Debug.Log(message);
            }
        }

        /// <summary>
        /// 打印error日志
        /// </summary>
        /// <param name="message"></param>
        public static void Error(string message)
        {
            if (Enabled)
            {
                Debug.LogError(message);
            }
        }

        /// <summary>
        /// 打印warning日志
        /// </summary>
        /// <param name="message"></param>
        public static void Warning(string message)
        {
            if (Enabled)
            {
                Debug.LogWarning(message);
            }
        }
    }
    
}