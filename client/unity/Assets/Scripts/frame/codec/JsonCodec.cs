using System;
using System.Text;
using UnityEngine;
using Object = System.Object;

namespace Nova.Codec
{
    /// <summary>
    /// 可能需要使用高级点的第三方json库，例如 Newtonsoft.Json，以支持复杂的对象
    /// </summary>
    public class JsonCodec : MessageCodec
    {
        public Object Decode(Type type, byte[] data)
        {
            object instance = Activator.CreateInstance(type);
            JsonUtility.FromJsonOverwrite(Encoding.UTF8.GetString(data), instance);
            return instance;
        }

        public byte[] Encode(Object data)
        {
            return Encoding.UTF8.GetBytes(JsonUtility.ToJson(data));
        }
    }
}