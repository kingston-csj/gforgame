using System.Collections.Generic;
using UnityEngine;

namespace Nova.Ui
{
    /// <summary>
    /// 面板打开时传送的参数类
    /// 用于传递任意参数
    /// </summary>
    public class PanelOpenIntend
    {
        private readonly Dictionary<string, object> _params = new Dictionary<string, object>();

        public void SetParam<T>(string key, T value)
        {
            if (string.IsNullOrEmpty(key))
            {
                Debug.LogError("Param key cannot be null or empty!");
                return;
            }

            _params[key] = value;
        }

        public T GetParam<T>(string key, T defaultValue = default)
        {
            if (!_params.ContainsKey(key)) return defaultValue;
            try
            {
                return (T)_params[key];
            }
            catch
            {
                Debug.LogError($"Failed to get param {key} with type {typeof(T)}");
                return defaultValue;
            }
        }

        public void Clear() => _params.Clear();
        
        public bool HasParam(string key) => _params.ContainsKey(key);
    }
}