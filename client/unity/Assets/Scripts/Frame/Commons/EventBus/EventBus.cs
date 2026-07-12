namespace Frame.Commons.EventBus
{
    using System;
    using System.Collections.Generic;
    using UnityEngine;

    /// <summary>
    /// 极简版全局事件总线
    /// 特性：string事件类型、无锁、无线程处理、全局唯一
    /// 适用：单主线程Unity项目，无需多线程/线程安全的场景
    /// </summary>
    public static class EventBus
    {
        #region 核心数据结构

        // 事件回调委托（统一传object类型参数，适配任意事件数据）
        public delegate void EventCallback(object eventData);

        // 事件映射表：Key=事件类型(string)，Value=该事件的所有回调
        private static readonly Dictionary<string, List<EventCallback>> _eventMap =
            new Dictionary<string, List<EventCallback>>();

        #endregion

        #region 公开API - 订阅事件

        /// <summary>
        /// 订阅事件
        /// </summary>
        /// <param name="eventType">事件类型（自定义string标识，如"PanelOpen"）</param>
        /// <param name="callback">回调方法（参数为object，可强转为自定义数据类型）</param>
        public static void Subscribe(string eventType, EventCallback callback)
        {
            if (string.IsNullOrEmpty(eventType) || callback == null)
            {
                Debug.LogError("EventBus: 事件类型或回调不能为空！");
                return;
            }

            // 初始化事件类型的回调列表
            if (!_eventMap.ContainsKey(eventType))
            {
                _eventMap[eventType] = new List<EventCallback>();
            }

            // 避免重复订阅
            if (_eventMap[eventType].Contains(callback))
            {
                Debug.LogWarning($"EventBus: 事件 {eventType} 已订阅该回调，无需重复订阅");
                return;
            }

            _eventMap[eventType].Add(callback);
        }

        #endregion

        #region 公开API - 取消订阅

        /// <summary>
        /// 取消指定事件的指定回调
        /// </summary>
        public static void Unsubscribe(string eventType, EventCallback callback)
        {
            if (string.IsNullOrEmpty(eventType) || callback == null) return;

            if (!_eventMap.ContainsKey(eventType)) return;

            // 移除回调
            if (_eventMap[eventType].Contains(callback))
            {
                _eventMap[eventType].Remove(callback);
                // 清空空列表，节省内存
                if (_eventMap[eventType].Count == 0)
                {
                    _eventMap.Remove(eventType);
                }
            }
        }

        /// <summary>
        /// 取消指定事件类型的所有回调
        /// </summary>
        public static void UnsubscribeAll(string eventType)
        {
            if (string.IsNullOrEmpty(eventType)) return;

            if (_eventMap.ContainsKey(eventType))
            {
                _eventMap.Remove(eventType);
            }
        }

        /// <summary>
        /// 取消所有事件的所有回调（慎用！）
        /// </summary>
        public static void UnsubscribeAll()
        {
            _eventMap.Clear();
        }

        #endregion

        #region 公开API - 发布事件

        /// <summary>
        /// 发布事件
        /// </summary>
        /// <param name="eventType">事件类型</param>
        /// <param name="eventData">事件数据（可传null/自定义类实例）</param>
        public static void Publish(string eventType, object eventData = null)
        {
            if (string.IsNullOrEmpty(eventType))
            {
                Debug.LogError("EventBus: 事件类型不能为空！");
                return;
            }

            if (!_eventMap.ContainsKey(eventType)) return;

            // 复制回调列表，避免遍历中列表被修改（即使无锁，也能避免基础异常）
            List<EventCallback> callbacks = new List<EventCallback>(_eventMap[eventType]);

            // 执行所有回调
            foreach (var callback in callbacks)
            {
                try
                {
                    callback.Invoke(eventData);
                }
                catch (Exception e)
                {
                    Debug.LogError($"EventBus: 执行事件 {eventType} 回调出错：{e.Message}\n{e.StackTrace}");
                }
            }
        }

        #endregion
    }
}