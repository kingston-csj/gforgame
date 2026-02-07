namespace Nova.Data
{
    /// <summary>
    /// 索引元数据接口：定义索引的名称、取值逻辑和唯一性
    /// </summary>
    /// <typeparam name="I">配置项类型（继承自AbsConfigItem）</typeparam>
    public interface IIndexMeta<I> where I : AbsConfigData
    {
        string Name { get; }
        object GetValue(I item);
        bool IsUnique { get; }
    }
}