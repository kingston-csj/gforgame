using System;
using UnityEngine;

namespace frame.Assets
{
    /// <summary>
    /// 通过unity直接配置的资源项
    /// </summary>
    /// <typeparam name="E"></typeparam>
    public class AssetResourceItem<E>
    {
        [Header("资源组")]
        public string group;
        
        [Header("资源描述")]
        public string desc;
        
        [Header("资源池")]
        public E[] pool;
    }
}