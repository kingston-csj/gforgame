using frame.Assets;
using Game.Configs;
using Nova.Net.Socket;

namespace Game.Core
{
    public class AppContext
    {
        //游戏引擎
        public static Engine engine;

        //资源工厂，启动时加载各种资源，如文本、图片、音频等
        public static AssetResourceFactory assetResourceFactory;

        // 配置数据管理器
        public static DataManager dataManager;
        
        // 游戏配置
        public static GameConfig gameConfig;
        
        /// <summary>
        /// 网络socket客户端
        /// </summary>
        public static SocketClient socketClient;
        
        
    }
}