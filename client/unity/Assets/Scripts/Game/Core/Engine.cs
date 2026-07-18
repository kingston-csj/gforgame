using System;
using System.Collections.Generic;
using System.IO;
using System.Reflection;
using frame.Assets;
using Game.Configs;
using Game.Login;
using Game.Net;
using Game.Net.Message;
using Nova.Codec;
using Nova.Commons.Util;
using Nova.Logger;
using Nova.Net.Socket;
using Nova.Ui;
using UnityEngine;

namespace Game.Core
{
    public class Engine : MonoBehaviour, IUiFactory
    {
        [Header("UI层列表")] public UiLayer[] layers;

        [Header("是否开启调试日志")]
        /// <summary>
        /// 是否开启日志
        /// </summary>
        public bool DebugLog = true;

        private void Awake()
        {
            AppContext.engine = this;
            // 资源工厂，启动时加载各种资源，如文本、图片、音频等
            AssetResourceFactory _assetResourceFactory = Resources.Load<AssetResourceFactory>("AssetResourceBinding");
            AppContext.assetResourceFactory = _assetResourceFactory;
            // 配置数据管理器
            DataManager dataManager = new DataManager();
            dataManager.AutoInit();
            AppContext.dataManager = dataManager;
            CommonValueMgr.Instance.AutoInject();

            // 游戏配置
            GameConfig gameConfig = new GameConfig();
            AppContext.gameConfig = gameConfig;

            string excelPath = Path.Combine(Application.dataPath, "Config/common.xlsx");
            if (!File.Exists(excelPath))
            {
                Debug.LogError("配置表不存在：" + excelPath);
                return;
            }

            List<PanelPrefabConfig> _panelConfigs = PanelModules.GetAllPanelConfigs();
            PanelManager.Init(this, _panelConfigs);

            // 异步打开登录面板 
            PanelManager.Instance.OpenPanel<LoginView>(R.Panels.Login);
        }
        public Transform GetNodeByLayer(LayerIds layer)
        {
            return layers[(int)layer].node;
        }
    }
}