using Nova.Data;
using Nova.Net.Socket;

namespace Game.Configs
{
    [DataTable(name ="common")]
    public class CommonData : AbsConfigData
    {
        public string key;

        public string value;

    }
}