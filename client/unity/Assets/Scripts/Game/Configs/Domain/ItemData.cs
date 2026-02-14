using Nova.Data;
using Nova.Net.Socket;

namespace Game.Configs
{
    [DataTable(name = "item")]
    /// <summary>
    /// 道具表
    /// </summary>
    public class ItemData : AbsConfigData
    {
        public int id;

        public string icon;

        public int type;

        public int quality;
    }
}