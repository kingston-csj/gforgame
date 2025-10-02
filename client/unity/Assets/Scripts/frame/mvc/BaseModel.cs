using System;
using System.Collections.Generic;

namespace Nova.Mvc
{
    public class BaseModel
    {
        // Register change listener for any field
        private Dictionary<string, LinkedList<Action<object>>> _changeCallbacks = new Dictionary<string, LinkedList<Action<object>>>();
        
        public void Register(string fieldName, Action<object> callback)
        {
            if (!_changeCallbacks.ContainsKey(fieldName))
            {
                _changeCallbacks[fieldName] = new LinkedList<Action<object>>();
            }
            _changeCallbacks[fieldName].AddLast(callback);
        }
        
        public void Notify(string fieldName, object value)
        {
            if (_changeCallbacks.ContainsKey(fieldName))
            {
                foreach (var callback in _changeCallbacks[fieldName])
                {
                    callback(value);
                }
            }
        }
    }
}