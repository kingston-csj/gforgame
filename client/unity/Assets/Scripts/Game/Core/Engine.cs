using System;
using System.Reflection;
using frame.Assets;
using Game.Configs;
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
    public class Engine : MonoBehaviour
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
            // 游戏配置
            GameConfig gameConfig = new GameConfig();
            AppContext.gameConfig = gameConfig;
            // 网络连接
            _CreateSocketClient();
        }

        private void _CreateSocketClient()
        {
            // 连接服务器
            SocketRuntimeEnvironment runtimeEnvironment =
                new SocketRuntimeEnvironment(typeof(MessageRouter), new JsonCodec(), new MessageFactory());
            // 自动注册所有的消息类型
            foreach (Type item in ClassScanner.ListAllSubclasses("Scripts\\Game\\Net", typeof(Message)))
            {
                // 获得class对应MessageMeta特性的cmd
                MessageMeta messageMeta = item.GetCustomAttribute(typeof(MessageMeta)) as MessageMeta;
                int cmd = messageMeta.Cmd;
                runtimeEnvironment.MessageFactory.Register(cmd, item);
            }
            SocketClient webSocketClient = new WebSocketClient(AppContext.gameConfig.serverUrl, runtimeEnvironment);
            AppContext.socketClient = webSocketClient;

            webSocketClient.ConnectAsync(() =>
            {
                // 连接成功
                LoggerUtil.Info("连接成功");

                // 发送登录请求
                ReqPlayerLogin reqLogin = new ReqPlayerLogin {playerId = "1001"};
                webSocketClient.Send(reqLogin, (ResPlayerLogin res) => Debug.Log($"登录成功，玩家名称：{res.name}"));
            });
        }
    }
}