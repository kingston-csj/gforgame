using System;
using frame.Assets;
using Game.Configs;
using Nova.Ui;
using UnityEngine;

namespace Game.Core
{
    public class Engine: MonoBehaviour
    {
        
        [Header("UI层列表")]
        public UiLayer[] layers;
        
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
        }
    }
}