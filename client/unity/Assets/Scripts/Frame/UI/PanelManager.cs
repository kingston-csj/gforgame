using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using UnityEngine;

namespace Nova.Ui
{
    // 面板管理器 
    public class PanelManager : IDisposable
    {
        // 单例实例（手动初始化）
        private static PanelManager _instance;

        public static PanelManager Instance
        {
            get
            {
                if (_instance == null)
                {
                    throw new InvalidOperationException(
                        "PanelManager has not been initialized! Call PanelManager.Init() first.");
                }

                return _instance;
            }
        }

        // 核心依赖
        private readonly IUiFactory _uiFactory; // UI工厂（用于动态获取层级节点）

        // 配置与缓存
        private readonly List<PanelPrefabConfig> _panelConfigs = new List<PanelPrefabConfig>();

        private readonly Dictionary<string, GameObject>
            _panelPrefabCache = new Dictionary<string, GameObject>(); // 面板名称 -> 预制体缓存

        private readonly Dictionary<string, IUiPanel> _panelCache = new Dictionary<string, IUiPanel>(); // 面板名称 -> 实例缓存

        private readonly Dictionary<string, Task<IUiPanel>> _loadingTasks = new Dictionary<string, Task<IUiPanel>>();


        /// <summary>
        /// 手动初始化面板管理器（必须先调用）
        /// </summary>
        /// <param name="uiFactory">UI工厂（用于获取层级节点）</param>
        /// <param name="panelConfigs">面板配置列表（路径版）</param>
        public static void Init( IUiFactory uiFactory,
            List<PanelPrefabConfig> panelConfigs = null)
        {
            if (_instance != null)
            {
                Debug.LogWarning("PanelManager has already been initialized!");
                return;
            }
            if (uiFactory == null)
            {
                throw new ArgumentNullException(nameof(uiFactory), "UI factory cannot be null!");
            }

            _instance = new PanelManager(uiFactory,
                panelConfigs ?? new List<PanelPrefabConfig>());
        }

        // 私有构造函数（单例）
        private PanelManager(IUiFactory uiFactory,
            List<PanelPrefabConfig> panelConfigs)
        {
            _uiFactory = uiFactory;
            _panelConfigs = panelConfigs;

            // 校验配置合法性
            ValidatePanelConfigs();
        }

        // 校验面板配置（路径不能为空/格式错误）
        private void ValidatePanelConfigs()
        {
            foreach (var config in _panelConfigs)
            {
                if (string.IsNullOrEmpty(config.PanelName))
                {
                    Debug.LogError("Invalid panel config: PanelName is null or empty!");
                    continue;
                }

                if (string.IsNullOrEmpty(config.PrefabPath))
                {
                    Debug.LogError($"Panel {config.PanelName} has empty PrefabPath!");
                    continue;
                }

                // 简单校验路径格式（避免包含后缀）
                if (config.PrefabPath.EndsWith(".prefab"))
                {
                    config.PrefabPath = config.PrefabPath.Replace(".prefab", "");
                    Debug.LogWarning(
                        $"Panel {config.PanelName} PrefabPath contains .prefab suffix, auto removed: {config.PrefabPath}");
                }

                // 校验层级节点是否存在
                var layerNode = _uiFactory.GetNodeByLayer(config.Layer);
                if (layerNode == null)
                {
                    Debug.LogError(
                        $"Panel {config.PanelName} layer {config.Layer} node is null! Check IUiFactory implementation.");
                }
            }
        }

        // 异步打开面板
        public async Task<T> OpenPanel<T>(string panelName, PanelOpenIntend openIntend = null, Action completeCallback = null)
            where T : class, IUiPanel
        {
            if (string.IsNullOrEmpty(panelName)) return null;

            // 防重复加载
            if (_loadingTasks.TryGetValue(panelName, out var task))
            {
                await task;
                return _panelCache.TryGetValue(panelName, out var cachedPanel) ? cachedPanel as T : null;
            }

            // 检查缓存
            if (_panelCache.TryGetValue(panelName, out var panel))
            {
                panel.Show(openIntend, completeCallback); // 调用IUIPanel的同步Show
                return panel as T;
            }

            // 加载面板
            Task<IUiPanel> loadTask = LoadPanelAndConvertAsync<T>(panelName, openIntend, completeCallback);
            _loadingTasks[panelName] = loadTask;

            try
            {
                // 等待任务完成后转换为目标类型
                var result = await loadTask;
                return result as T;
            }
            finally
            {
                _loadingTasks.Remove(panelName);
            }
        }

        // 新增包装方法：将 Task<T> 转换为 Task<IUIPanel>
        private async Task<IUiPanel> LoadPanelAndConvertAsync<T>(string panelName, PanelOpenIntend openIntend, Action completeCallback = null)
            where T : class, IUiPanel
        {
            // 调用原有加载逻辑，返回 T 类型
            var panel = await LoadPanelAsync<T>(panelName, openIntend, completeCallback);
            // 向上转换为 IUIPanel（安全，因为 T 约束为 IUIPanel）
            return panel as IUiPanel;
        }

        // 异步加载面板（内部）
        private async Task<T> LoadPanelAsync<T>(string panelName, PanelOpenIntend openIntend, Action completeCallback = null)
            where T : class, IUiPanel
        {
            // 1. 获取面板配置
            var config = _panelConfigs.Find(c => c.PanelName == panelName);
            if (config == null)
            {
                Debug.LogError($"Panel config not found for {panelName}!");
                return null;
            }

            // 2. 获取面板所属层级的父节点
            var parentNode = _uiFactory.GetNodeByLayer(config.Layer);
            if (parentNode == null)
            {
                Debug.LogError($"Cannot load panel {panelName}: layer {config.Layer} node is null!");
                return null;
            }

            // 3. 异步加载预制体
            GameObject prefab = await LoadPrefabAsync(config.PrefabPath, panelName);
            if (prefab == null) return null;
            GameObject panelObj = GameObject.Instantiate(prefab, parentNode, false);

            // 4. 获取IUIPanel组件
            T panel = panelObj.GetComponent<T>();
            if (panel == null)
            {
                Debug.LogError($"Panel {panelName} does not implement {typeof(T).Name} (IUIPanel)!");
                return null;
            }

            // 6. 初始化并显示面板
            panel.Init();
            panel.Show(openIntend, completeCallback);

            // 7. 加入缓存
            _panelCache[panelName] = panel;

            return panel;
        }

        // 异步加载预制体
        private async Task<GameObject> LoadPrefabAsync(string path, string panelName)
        {
            if (_panelPrefabCache.TryGetValue(panelName, out var prefab))
                return prefab;

            // 步骤1：在后台启动异步加载（不阻塞主线程）
            ResourceRequest request = Resources.LoadAsync<GameObject>(path);
            if (request == null)
            {
                Debug.LogError($"启动加载失败：路径={path}（面板={panelName}）");
                return null;
            }

            // 步骤2：异步等待加载完成（真正的await，不阻塞主线程）
            while (!request.isDone)
            {
                await Task.Yield(); // 此时Yield有效，因为不在同步闭包内
            }

            // 步骤3：在主线程获取加载结果
            prefab = request.asset as GameObject;
            if (prefab == null)
            {
                Debug.LogError(
                    $"加载预制体失败：路径={path}（面板={panelName}），请检查：1.路径是否正确 2.预制体是否在Resources文件夹下 3.路径是否不含.prefab后缀");
                return null;
            }

            _panelPrefabCache[panelName] = prefab;
            return prefab;
        }

        // 关闭面板
        public void ClosePanel(string panelName)
        {
            if (string.IsNullOrEmpty(panelName)) return;

            if (_panelCache.TryGetValue(panelName, out var panel))
                panel.Hide();
            else
                Debug.LogWarning($"Panel {panelName} not found in cache, cannot close!");
        }

        // 销毁面板
        public void DestroyPanel(string panelName)
        {
            if (string.IsNullOrEmpty(panelName)) return;

            if (_panelCache.TryGetValue(panelName, out var panel))
            {
                panel.DestroyPanel();
                _panelCache.Remove(panelName);
            }
            else
                Debug.LogWarning($"Panel {panelName} not found in cache, cannot destroy!");
        }

        /// <summary>
        /// 清空所有面板
        /// </summary>
        public void ClearAllPanels()
        {
            foreach (var panel in _panelCache.Values)
                panel.DestroyPanel();

            _panelCache.Clear();
            _loadingTasks.Clear();
        }

        public void Dispose()
        {
            ClearAllPanels();
            _panelPrefabCache.Clear();
            _panelConfigs.Clear();
            _instance = null;

            Debug.Log("PanelManager disposed successfully!");
        }
    }
}