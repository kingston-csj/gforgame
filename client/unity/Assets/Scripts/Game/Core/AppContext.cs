using frame.Assets;
using Game.Configs;

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
    }
}