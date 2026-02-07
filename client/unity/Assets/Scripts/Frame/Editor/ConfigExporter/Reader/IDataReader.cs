using System;
using System.Collections.Generic;
using System.IO;

namespace Nova.Editor.ConfigExporter
{
    public interface IDataReader
    {
        /// <summary>
        /// 将文件流转换为指定类型的配置记录集合
        /// </summary>
        /// <typeparam name="T">配置实体类类型</typeparam>
        /// <param name="stream">文件流（Excel文件流）</param>
        /// <param name="type">配置实体类Type</param>
        /// <returns>配置记录集合</returns>
        List<T> Read<T>(Stream stream, Type type) where T : new();
    }
}