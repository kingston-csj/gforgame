using System;
using System.Collections.Generic;
using UnityEngine;

namespace Nova.Mvc
{
    public class BaseModel
{
    // 字段变更回调字典
    private Dictionary<string, LinkedList<Action<object>>> _changeCallbacks = new Dictionary<string, LinkedList<Action<object>>>();

    // 注册字段变更回调
    public void Register(string fieldName, Action<object> callback)
    {
        if (string.IsNullOrEmpty(fieldName))
        {
            Debug.LogError("Field name cannot be null or empty!");
            return;
        }
        
        if (callback == null)
        {
            Debug.LogError("Callback cannot be null!");
            return;
        }

        if (!_changeCallbacks.ContainsKey(fieldName))
        {
            _changeCallbacks[fieldName] = new LinkedList<Action<object>>();
        }
        
        // 避免重复注册
        if (!_changeCallbacks[fieldName].Contains(callback))
        {
            _changeCallbacks[fieldName].AddLast(callback);
        }
    }

    // 移除字段变更回调
    public void Unregister(string fieldName, Action<object> callback)
    {
        if (string.IsNullOrEmpty(fieldName) || callback == null)
        {
            return;
        }

        if (_changeCallbacks.ContainsKey(fieldName))
        {
            _changeCallbacks[fieldName].Remove(callback);
            // 如果回调列表为空，移除该字段
            if (_changeCallbacks[fieldName].Count == 0)
            {
                _changeCallbacks.Remove(fieldName);
            }
        }
    }

    // 通知字段变更
    public void Notify(string fieldName, object value)
    {
        if (string.IsNullOrEmpty(fieldName))
        {
            Debug.LogError("Field name cannot be null or empty!");
            return;
        }

        if (_changeCallbacks.TryGetValue(fieldName, out var callbacks))
        {
            // 遍历副本避免回调中修改列表导致的异常
            var callbackList = new List<Action<object>>(callbacks);
            foreach (var callback in callbackList)
            {
                try
                {
                    callback?.Invoke(value);
                }
                catch (Exception e)
                {
                    Debug.LogError($"Callback execution failed for field {fieldName}: {e.Message}");
                }
            }
        }
    }

    // 清空所有回调
    public virtual void ClearCallbacks()
    {
        _changeCallbacks.Clear();
    }
}

}