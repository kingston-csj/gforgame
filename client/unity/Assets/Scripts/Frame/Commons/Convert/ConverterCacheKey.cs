using System;

namespace Nova.Commons.Convert
{
    public class ConverterCacheKey
    {
        public Type SourceType;
        public Type TargetType;

        public override int GetHashCode()
        {
            // 方案1：C# 7.2+ 推荐（Unity 2018+支持）
            return HashCode.Combine(SourceType, TargetType);

            // 方案2：兼容旧版本Unity
            // int hash = 17;
            // hash = hash * 31 + (SourceType?.GetHashCode() ?? 0);
            // hash = hash * 31 + (TargetType?.GetHashCode() ?? 0);
            // return hash;
        }

        // 重写Equals：比较SourceType和TargetType的值
        public override bool Equals(object obj)
        {
            // 1. 引用相同，直接返回true
            if (ReferenceEquals(this, obj)) return true;

            // 2. obj为空或类型不同，返回false
            if (obj == null || GetType() != obj.GetType()) return false;

            // 3. 比较核心字段
            var other = (ConverterCacheKey)obj;
            return SourceType == other.SourceType && TargetType == other.TargetType;
        }

        // 可选：实现IEquatable<T>，提升性能（避免装箱）
        public bool Equals(ConverterCacheKey other)
        {
            if (other == null) return false;
            return SourceType == other.SourceType && TargetType == other.TargetType;
        }
    }
}