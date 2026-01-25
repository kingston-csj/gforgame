using System;
using System.Reflection;

namespace Nova.Net.Socket
{
    public class DefaultSocketIoDispatcher : SocketIoDispatcher
    {
        private SocketRuntimeEnvironment _socketRuntimeEnvironment;

        /// <summary>
        ///     构造函数
        /// </summary>
        /// <param name="runtimeEnvironment"> 指定socket运行时环境 </param>
        public DefaultSocketIoDispatcher(SocketRuntimeEnvironment runtimeEnvironment)
        {
            this._socketRuntimeEnvironment = runtimeEnvironment;
            // 自动扫描处理服务端主动推送的消息处理器
            _AutoCreateResponseDelegate(runtimeEnvironment.MessageRouterType);
        }

        /// <summary>
        ///     自动注册响应处理器
        /// </summary>
        private void _AutoCreateResponseDelegate(Type messageRouterType)
        {
            object router = Activator.CreateInstance(messageRouterType);
            Type currentType = router.GetType();
            MethodInfo[] methods = currentType.GetMethods(BindingFlags.Public | BindingFlags.Instance);
            MethodInfo methodInfo =
                currentType.GetMethod("RegisterCallbackDelegate", BindingFlags.Public | BindingFlags.Instance);
            foreach (MethodInfo method in methods)
            {
                // 检查是否有 MessageHandler 特性
                MessageHandler attr = method.GetCustomAttribute<MessageHandler>();
                if (attr != null)
                {
                    ParameterInfo[] parameters = method.GetParameters();
                    if (parameters.Length != 1)
                    {
                        throw new Exception($"自动注册推送处理器:{method.Name} 参数错误");
                    }

                    Type responseType = parameters[0].ParameterType;
                    // 获得Message类的实例方法GetCmd()返回值
                    int responseCmd = ((MessageMeta)responseType.GetCustomAttribute(typeof(MessageMeta))).Cmd;
                    // 绑定消息id与消息类型
                    _socketRuntimeEnvironment.MessageFactory.Register(responseCmd, responseType);
                    // 创建委托
                    Type delegateType = typeof(Action<>).MakeGenericType(responseType);
                    Delegate handler = Delegate.CreateDelegate(delegateType, router, method);

                    MethodInfo genericMethod = methodInfo.MakeGenericMethod(responseType);
                    genericMethod.Invoke(router, new object[] { responseCmd, handler });
                }
            }
        }

        public void OnOpen()
        {
        }

        public void OnClose()
        {
            throw new NotImplementedException();
        }

        public void OnMessage(SocketDataFrame dataFrame)
        {
            // 客户端请求的响应
            MessageCallback callback = CallbackMgr.Fetch(dataFrame.index);
            if (callback != null)
            {
                callback.callback(dataFrame.message);
            }
            else
            {
                // 服务器主动推送的消息
                MessageDispatcher.Dispatch(dataFrame.cmd, dataFrame.message);
            }
        }
    }
}