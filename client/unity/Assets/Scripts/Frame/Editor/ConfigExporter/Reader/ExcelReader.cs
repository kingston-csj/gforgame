using Nova.Commons.Convert;

namespace Nova.Editor.ConfigExporter
{
    using System;
    using System.Collections.Generic;
    using System.IO;
    using System.Reflection;
    using ExcelDataReader; // 引入ExcelDataReader库
    using UnityEngine;

    /// <summary>
    /// 基于ExcelDataReader的实现类
    /// </summary>
    public class ExcelReader : BaseDataReader, IDataReader
    {
        // 静态初始化：设置编码（解决中文乱码）
        static ExcelReader()
        {
            System.Text.Encoding.RegisterProvider(System.Text.CodePagesEncodingProvider.Instance);
        }

        public ExcelReader(IConversionService conversionService) : base(conversionService)
        {
        }

        // 实现Read方法
        public List<T> Read<T>(Stream stream, Type type) where T : new()
        {
            try
            {
                // 用ExcelDataReader打开流（兼容所有.xlsx/.xls）
                using (var reader = ExcelReaderFactory.CreateReader(stream))
                {
                    // 转换为DataSet（方便按行/列读取）
                    var result = reader.AsDataSet(new ExcelDataSetConfiguration
                    {
                        ConfigureDataTable = (_) => new ExcelDataTableConfiguration
                        {
                            UseHeaderRow = false, // 自己处理表头，不使用库的表头逻辑
                            ReadHeaderRow = (rowReader) => { }
                        }
                    });

                    // 校验工作表
                    if (result.Tables.Count == 0)
                        throw new Exception("Excel中无有效工作表");

                    var table = result.Tables[0]; // 取第一个工作表
                    int lastRow = table.Rows.Count;
                    int lastCol = table.Columns.Count;
                    if (lastRow < 1 || lastCol < 1)
                        throw new Exception("Excel中无有效数据");

                    bool hasHeader = false;
                    CellHeader[] headers = null;
                    List<CellColumn[]> records = new List<CellColumn[]>();
                    string[] exportTypes = Array.Empty<string>();

                    // 遍历所有行（ExcelDataReader行索引从0开始）
                    for (int rowIdx = 0; rowIdx < lastRow; rowIdx++)
                    {
                        // 跳过空行
                        if (IsRowEmpty(table, rowIdx, lastCol))
                            continue;

                        // 获取第一列标记值（统一转小写）
                        string firstCell = GetCellValue(table, rowIdx, 0)?.Trim().ToLower() ?? "";

                        // 解析表头行（header标记）
                        if (firstCell == BEGIN){
                            headers = ParseHeader(table, rowIdx, lastCol, type);
                            hasHeader = true;
                            continue;
                        }

                        // 解析导出类型行（export标记）
                        if (firstCell == EXPORT)
                        {
                            exportTypes = ParseExportTypes(table, rowIdx, lastCol, headers?.Length ?? 0);
                            continue;
                        }

                        // 表头未解析完成则跳过
                        if (!hasHeader || headers == null)
                            continue;

                        // 遇到结束标记则停止（可选）
                        if (firstCell == END)
                            break;

                        // 解析数据行
                        CellColumn[] rowData = ParseDataRow(table, rowIdx, lastCol, headers, exportTypes);
                        if (rowData != null)
                            records.Add(rowData);
                    }

                    // 复用基类的ReadRecords转换为实体列表
                    return ReadRecords<T>(type, exportTypes, records, 0);
                }
            }
            catch (Exception ex)
            {
                throw new Exception($"Excel解析失败：{ex.Message}", ex);
            }
        }

        #region 私有辅助方法（适配ExcelDataReader的DataTable）

        /// <summary>
        /// 解析表头（适配DataTable）
        /// </summary>
        private CellHeader[] ParseHeader(System.Data.DataTable table, int rowIdx, int lastCol, Type entityType)
        {
            List<CellHeader> headers = new List<CellHeader>();
            // 从第1列开始（第0列是标记列）
            for (int colIdx = 1; colIdx < lastCol; colIdx++)
            {
                string colName = GetCellValue(table, rowIdx, colIdx)?.Trim() ?? "";
                if (string.IsNullOrEmpty(colName))
                    continue;

                // 复用基类的字段查找逻辑
                FieldInfo field = FindFieldInClassHierarchy(entityType, colName);
                if (field == null && !IgnoreUnknownFields)
                    throw new Exception($"表头字段不存在：{colName}（实体：{entityType.Name}）");

                headers.Add(new CellHeader
                {
                    Column = colName,
                    Field = field
                });
            }

            return headers.ToArray();
        }

        /// <summary>
        /// 解析导出类型（适配DataTable）
        /// </summary>
        private string[] ParseExportTypes(System.Data.DataTable table, int rowIdx, int lastCol, int headerCount)
        {
            List<string> exportTypes = new List<string>();
            // 从第1列开始
            for (int colIdx = 1; colIdx < lastCol; colIdx++)
            {
                string type = GetCellValue(table, rowIdx, colIdx)?.Trim().ToLower() ?? EXPORT_TYPE_NONE;
                exportTypes.Add(type);
            }

            // 补全导出类型（与表头数量一致）
            while (exportTypes.Count < headerCount)
                exportTypes.Add(EXPORT_TYPE_NONE);

            return exportTypes.ToArray();
        }

        /// <summary>
        /// 解析数据行（适配DataTable）
        /// </summary>
        private CellColumn[] ParseDataRow(System.Data.DataTable table, int rowIdx, int lastCol, CellHeader[] headers,
            string[] exportTypes)
        {
            List<CellColumn> rowData = new List<CellColumn>();
            for (int i = 0; i < headers.Length; i++)
            {
                int colIdx = i + 1; // 第1列对应第一个表头字段
                if (colIdx >= lastCol)
                {
                    rowData.Add(null);
                    continue;
                }

                // 跳过无需导出的字段
                if (i < exportTypes.Length && exportTypes[i] != EXPORT_TYPE_BOTH &&
                    exportTypes[i] != EXPORT_TYPE_SERVER)
                {
                    rowData.Add(null);
                    continue;
                }

                // 读取单元格值
                string cellValue = GetCellValue(table, rowIdx, colIdx)?.Trim() ?? "";
                rowData.Add(new CellColumn
                {
                    Header = headers[i],
                    Value = cellValue
                });
            }

            return rowData.ToArray();
        }

        /// <summary>
        /// 获取DataTable单元格值（兜底空值）
        /// </summary>
        private string GetCellValue(System.Data.DataTable table, int rowIdx, int colIdx)
        {
            try
            {
                var cellValue = table.Rows[rowIdx][colIdx];
                return cellValue == DBNull.Value ? "" : cellValue.ToString();
            }
            catch
            {
                return "";
            }
        }

        /// <summary>
        /// 判断DataTable行是否为空
        /// </summary>
        private bool IsRowEmpty(System.Data.DataTable table, int rowIdx, int lastCol)
        {
            for (int colIdx = 0; colIdx < lastCol; colIdx++)
            {
                if (!string.IsNullOrEmpty(GetCellValue(table, rowIdx, colIdx)))
                    return false;
            }

            return true;
        }

        #endregion
    }
}