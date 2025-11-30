using System;
using AppContext = Game.Core.AppContext;

/// <summary>
/// 对common表配置的自动注入
/// 对于该类的所有字段，会自动从common表配置中注入相同key的值
/// </summary>
public class CommonValueMgr
{
    public static CommonValueMgr Instance = new ();

    /// <summary>
    /// 招募令不足消耗的等价钻石
    /// </summary>
    public int heroRecruitDiamond;

    public void AutoInject()
    {
        // 遍历所有字段
        foreach (var field in GetType().GetFields())
        {
            if (field.IsStatic) continue;
            // 从common表中获取相同key的值
            var config = AppContext.dataManager.configCommonContainer.GetValue(field.Name);
            if (config != null)
            {
                // 注入值
                if (field.FieldType == typeof(int))
                {
                    field.SetValue(this, Convert.ToInt32(config));
                }
                else if (field.FieldType == typeof(string))
                {
                    field.SetValue(this, Convert.ToString(config));
                }
                else if (field.FieldType == typeof(bool))
                {
                    field.SetValue(this, Convert.ToBoolean(config));
                }
                else if (field.FieldType == typeof(float))
                {
                    field.SetValue(this, Convert.ToSingle(config));
                }
                else
                {
                    throw new Exception($"字段 {field.Name} 的类型 {field.GetType()} 不支持注入！");
                }
            }
        }
    }
}