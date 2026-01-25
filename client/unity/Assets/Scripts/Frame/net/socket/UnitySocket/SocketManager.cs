using System;
using System.Collections.Generic;
using UnityEngine;

namespace Nova.Net.UnitySocket
{
    public class SocketManager : MonoBehaviour
    {
        private const string RootName = "[UnitySocket]";
        private static SocketManager _instance;
        private readonly List<UnitySocket> _sockets = new();
        private readonly object _socketListLock = new(); // 列表操作锁

        public static SocketManager Instance
        {
            get
            {
                if (_instance == null)
                    CreateInstance();
                return _instance;
            }
        }

        private void Awake()
        {
            if (_instance != null && _instance != this)
            {
                Destroy(gameObject);
                return;
            }

            _instance = this;
            DontDestroyOnLoad(gameObject);
        }

        public static void CreateInstance()
        {
            GameObject go = GameObject.Find("/" + RootName);
            if (go == null)
                go = new GameObject(RootName);

            _instance = go.GetComponent<SocketManager>();
            if (_instance == null)
                _instance = go.AddComponent<SocketManager>();
        }

        /// <summary>
        /// 线程安全添加Socket
        /// </summary>
        public void Add(UnitySocket unitySocket)
        {
            if (unitySocket == null) return;

            lock (_socketListLock)
            {
                if (!_sockets.Contains(unitySocket))
                    _sockets.Add(unitySocket);
            }
        }

        /// <summary>
        /// 线程安全移除Socket
        /// </summary>
        public void Remove(UnitySocket unitySocket)
        {
            if (unitySocket == null) return;

            lock (_socketListLock)
            {
                if (_sockets.Contains(unitySocket))
                    _sockets.Remove(unitySocket);
            }
        }

        private void Update()
        {
            // 复制一份列表，避免遍历中修改原列表
            UnitySocket[] socketsCopy;
            lock (_socketListLock)
            {
                socketsCopy = _sockets.ToArray();
            }

            foreach (var socket in socketsCopy)
            {
                try
                {
                    socket.Update();
                }
                catch (Exception ex)
                {
                    Debug.LogError($"Socket更新异常：{ex.Message}");
                }
            }
        }

        // 场景切换/销毁时清理所有Socket
        private void OnDestroy()
        {
            lock (_socketListLock)
            {
                foreach (var socket in _sockets)
                {
                    socket?.Dispose();
                }

                _sockets.Clear();
            }
        }
    }
}