using Nova.Commons.Convert;

namespace Nova.Editor.ConfigExporter
{
    using System;
    using System.Collections.Generic;
    using System.Reflection;
    using System.Text;

    /// <summary>
    /// 配置数据读取基类
    /// </summary>
    public abstract class BaseDataReader
    {
        #region 常量定义

        /// <summary>
        /// 表头开始标记（Excel中第一列值为HEADER时，该行是字段名列）
        /// </summary>
        protected const string BEGIN = "header";

        /// <summary>
        /// 配置表结束标记（Excel中第一列值为END时，停止读取）
        /// </summary>
        protected const string END = "end";

        /// <summary>
        /// 导出类型行标记（Excel中第一列值为EXPORT时，该行是导出类型列）
        /// </summary>
        protected const string EXPORT = "export";

        /// <summary>
        /// 导出类型：服务端+客户端均导入
        /// </summary>
        protected const string EXPORT_TYPE_BOTH = "both";

        /// <summary>
        /// 导出类型：仅服务端
        /// </summary>
        protected const string EXPORT_TYPE_SERVER = "server";

        /// <summary>
        /// 导出类型：仅客户端
        /// </summary>
        protected const string EXPORT_TYPE_CLIENT = "client";

        /// <summary>
        /// 导出类型：均不导入
        /// </summary>
        protected const string EXPORT_TYPE_NONE = "none";

        #endregion

        #region 字段

        /// <summary>
        /// 是否忽略无法识别的字段（配置表有但实体类无的字段）
        /// </summary>
        protected bool IgnoreUnknownFields = true;

        /// <summary>
        /// 类型转换服务
        /// </summary>
        protected readonly IConversionService dataConversionService;

        #endregion

        #region 构造函数

        public BaseDataReader(IConversionService conversionService)
        {
            dataConversionService = conversionService;
        }

        #endregion

        #region 核心方法

        /// <summary>
        /// 递归查找类及其父类的字段（对应Java的findFieldInClassHierarchy）
        /// </summary>
        /// <param name="type">目标类型</param>
        /// <param name="fieldName">字段名</param>
        /// <returns>字段信息</returns>
        /// <exception cref="MissingFieldException">未找到字段时抛出</exception>
        protected FieldInfo FindFieldInClassHierarchy(Type type, string fieldName)
        {
            if (type == null || type == typeof(object))
            {
                throw new MissingFieldException($"类型 {type?.FullName} 及其父类中未找到字段 {fieldName}");
            }

            // 查找当前类的私有/公有字段（忽略大小写，适配Excel可能的大小写不一致）
            var field = type.GetField(fieldName,
                BindingFlags.Public | BindingFlags.NonPublic | BindingFlags.Instance | BindingFlags.IgnoreCase);
            if (field != null)
            {
                return field;
            }

            // 递归查找父类
            return FindFieldInClassHierarchy(type.BaseType, fieldName);
        }

        /// <summary>
        /// 获取指定索引的导出类型
        /// </summary>
        /// <param name="exportHeaders">导出类型数组</param>
        /// <param name="index">列索引</param>
        /// <returns>导出类型</returns>
        protected string GetExportType(string[] exportHeaders, int index)
        {
            if (exportHeaders == null || exportHeaders.Length == 0 || index >= exportHeaders.Length)
            {
                return EXPORT_TYPE_BOTH; // 默认全导出
            }

            return exportHeaders[index]?.ToLower() ?? EXPORT_TYPE_BOTH;
        }

        /// <summary>
        /// 将Excel行数据转换为实体列表（对应Java的readRecords）
        /// </summary>
        /// <typeparam name="T">实体类型</typeparam>
        /// <param name="type">实体Type</param>
        /// <param name="exportHeaders">导出类型数组</param>
        /// <param name="rows">Excel行数据</param>
        /// <param name="headerIndex">表头行索引</param>
        /// <returns>实体列表</returns>
        protected List<T> ReadRecords<T>(Type type, string[] exportHeaders, List<CellColumn[]> rows, int headerIndex)
            where T : new()
        {
            var records = new List<T>(rows.Count);

            for (int i = 0; i < rows.Count; i++)
            {
                var record = rows[i];
                T obj = new T(); // C# 要求实体有无参构造函数

                for (int j = 0; j < record.Length; j++)
                {
                    var column = record[j];
                    if (column == null || string.IsNullOrEmpty(column.Header?.Column))
                    {
                        continue;
                    }

                    string colName = column.Header.Column;
                    string exportType = GetExportType(exportHeaders, j);

                    // 仅处理服务端/全导出的字段
                    if (exportType == EXPORT_TYPE_BOTH || exportType == EXPORT_TYPE_CLIENT)
                    {
                        try
                        {
                            FieldInfo field = FindFieldInClassHierarchy(type, colName);
                            if (field == null) continue;

                            // 设置字段可访问（私有字段也能赋值）
                            field.SetValue(obj, dataConversionService.Convert(column.Value, field.FieldType));
                        }
                        catch (MissingFieldException ex)
                        {
                            if (!IgnoreUnknownFields)
                            {
                                throw new Exception($"配置表 {type.Name} 第 {i + headerIndex + 2} 行字段 {colName} 不存在", ex);
                            }

                            Console.WriteLine($"警告：配置表 {type.Name} 第 {i + headerIndex + 2} 行字段 {colName} 不存在，已忽略");
                        }
                        catch (Exception ex)
                        {
                            Console.WriteLine(
                                $"错误：配置表 {type.Name} 第 {i + headerIndex + 2} 行字段 {colName} 转换异常：{ex.Message}");
                            throw;
                        }
                    }
                }

                records.Add(obj);
            }

            return records;
        }

        #endregion

        #region 公共方法

        /// <summary>
        /// 设置是否忽略未知字段
        /// </summary>
        /// <param name="ignoreUnknownFields">是否忽略</param>
        public void SetIgnoreUnknownFields(bool ignoreUnknownFields)
        {
            IgnoreUnknownFields = ignoreUnknownFields;
        }

        #endregion
    }
}