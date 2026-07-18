namespace Nova.Ui
{
    using System;
    using System.Collections.Generic;
    using System.Threading.Tasks;
    using UnityEngine;

    // Unity主线程调度器
    public class UnityMainThreadDispatcher : MonoBehaviour, IMainThreadDispatcher
    {
        // 单例实例
        private static UnityMainThreadDispatcher _instance;

        public static UnityMainThreadDispatcher Instance
        {
            get
            {
                if (_instance == null)
                {
                    var obj = new GameObject("UnityMainThreadDispatcher");
                    _instance = obj.AddComponent<UnityMainThreadDispatcher>();
                    DontDestroyOnLoad(obj);
                }

                return _instance;
            }
        }

        // 任务队列（线程安全）
        private readonly Queue<Action> _actionQueue = new Queue<Action>();
        private readonly object _queueLock = new object();

        private void Awake()
        {
            if (_instance != null && _instance != this)
            {
                Destroy(gameObject);
                return;
            }

            _instance = this;
        }

        private void Update()
        {
            // 每帧在主线程执行队列中的所有任务
            lock (_queueLock)
            {
                while (_actionQueue.Count > 0)
                {
                    try
                    {
                        _actionQueue.Dequeue().Invoke();
                    }
                    catch (Exception e)
                    {
                        Debug.LogError($"Execute main thread task failed: {e.Message}");
                    }
                }
            }
        }

        #region 实现IMainThreadDispatcher接口的所有方法

        /// <summary>
        /// 投递无返回值的操作到主线程执行
        /// </summary>
        public Task EnqueueAsync(Action action)
        {
            if (action == null)
            {
                return Task.CompletedTask;
            }

            // 创建任务完成源，用于通知调用方任务完成
            var tcs = new TaskCompletionSource<bool>();

            lock (_queueLock)
            {
                _actionQueue.Enqueue(() =>
                {
                    try
                    {
                        action.Invoke(); // 主线程执行任务
                        tcs.SetResult(true); // 标记任务成功完成
                    }
                    catch (Exception e)
                    {
                        tcs.SetException(e); // 传递异常给调用方
                    }
                });
            }

            return tcs.Task; // 返回Task，支持await等待
        }

        /// <summary>
        /// 投递有返回值的操作到主线程执行（补全的核心方法）
        /// </summary>
        public Task<T> EnqueueAsync<T>(Func<T> func)
        {
            if (func == null)
            {
                // 返回默认值的已完成任务
                return Task.FromResult(default(T));
            }

            // 创建带返回值的任务完成源
            var tcs = new TaskCompletionSource<T>();

            lock (_queueLock)
            {
                _actionQueue.Enqueue(() =>
                {
                    try
                    {
                        T result = func.Invoke(); // 主线程执行并获取返回值
                        tcs.SetResult(result); // 把返回值传递给调用方
                    }
                    catch (Exception e)
                    {
                        tcs.SetException(e); // 传递异常
                    }
                });
            }

            return tcs.Task; // 返回带返回值的Task
        }

        #endregion

        private void OnDestroy()
        {
            if (_instance == this)
            {
                _instance = null;
            }

            // 清空任务队列
            lock (_queueLock)
            {
                _actionQueue.Clear();
            }
        }
    }
}