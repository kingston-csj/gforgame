using System.Collections.Generic;
using Nova.Data;
using UnityEngine;

namespace Game.Configs
{
    public class ConfigCommonContainer : ConfigContainer<CommonData>
    {
        // common.json 的映射表
        private Dictionary<string, string> mapper = new();

        public ConfigCommonContainer(TextAsset textAsset) : base(textAsset)
        {
        }

        protected override void AfterLoad()
        {
            foreach (var item in GetItems())
            {
                mapper.Add(item.key, item.value);
            }
        }

        public string GetValue(string key)
        {
            return mapper[key];
        }
    }
}