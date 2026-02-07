using Nova.Data;
using UnityEngine;

namespace Game.Configs
{
    public class ConfigItemContainer:ConfigContainer<ItemData>
    {
        
        public ConfigItemContainer(TextAsset textAsset):base(textAsset)
        {
            
        }
    }
}