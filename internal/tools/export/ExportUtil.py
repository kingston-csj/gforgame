import glob
import os

import pandas as pd

import json

ignoreColumns = {"id"}

source_directory = '../../data/'
target_json_directory = "../../client/cocos/assets/resources/config"
target_typescript_directory = "../../client/cocos/assets/scripts/data/config/model"


class ExportUtil:

    @staticmethod
    def excel2json(excel_path):
        file_name_without_ext = os.path.splitext(os.path.basename(excel_path))[0]
        df = pd.read_excel(excel_path)
        type_row = df.iloc[0]  # 获取第二行（类型说明行）
        header_row = df.iloc[1]  # 获取第三行（属性名行）

        result = []
        for index, row in df[2:].iterrows():  # 从第四行开始处理实际数据
            item = {}
            try:
                for col_idx, col in enumerate(df.columns):
                    if col_idx == 0:
                        continue
                    if pd.isna(type_row[col]):
                        continue
                    col_type = type_row[col]
                    col_name = header_row[col]
                    value = row[col]
                    if 'int' in str(col_type):
                        item[col_name] = int(value) if pd.notnull(value) else 0
                    elif 'float' in str(col_type):
                        item[col_name] = float(value) if pd.notnull(value) else 0
                    elif 'str' in str(col_type) or 'string' in str(col_type):
                        item[col_name] = value if pd.notnull(value) else None
                    elif 'list' in str(col_type):
                        try:
                            item[col_name] = json.loads(value) if pd.notnull(value) else None
                        except json.JSONDecodeError:
                            item[col_name] = value
                    elif 'json' in str(col_type):
                        try:
                            item[col_name] = json.loads(value) if pd.notnull(value) else None
                        except json.JSONDecodeError:
                            item[col_name] = value
                    else:
                        item[col_name] = value
            except Exception as e:
                print(e)
                print(f"解析{excel_path}出错")
            result.append(item)

        file_name = file_name_without_ext[0].upper() + file_name_without_ext[1:]
        ExportUtil.json2TsFile(type_row, header_row, result[0], file_name_without_ext,
                               os.path.join(target_typescript_directory, f"{file_name}Data.ts"))

        return result

    @staticmethod
    def type2TypeScript(type):
        try:
            if "int" in type:
                return "number"
            elif "float" in type:
                return "number"
            elif "str" in type:
                return "string"
            elif "bool" in type:
                return "boolean"
            return "json"
        except Exception as e:
            return ""

    @staticmethod
    def json2TsFile(typeRow, headerRow, dataRow, class_name, target_code_path):
        """
        根据JSON数据生成对应ts代码文件
        """

        # 生成主配置类代码
        config_class_code = f'''
       import BaseConfigItem from '../BaseConfigItem';
            '''

        field_lists = []
        firstColumn = False
        for key, value in headerRow.items():
            if not (firstColumn):
                firstColumn = True
                continue
            if pd.isnull(typeRow[key]):
                continue
            if headerRow[key] in ignoreColumns:
                continue
            field_type = typeRow[key]
            ts_type = ExportUtil.type2TypeScript(field_type)
            if ts_type == "json":
                innerClassName = headerRow[key] + "Def"
                innerClassName = innerClassName[0].upper() + innerClassName[1:]
                config_class_code += f"\nexport class {innerClassName} {{\n"
                try:
                    for _key, _value in dataRow[headerRow[key]][0].items():
                        value_type = ExportUtil.type2TypeScript(type(_value).__name__)
                        config_class_code += f"    public {_key}: {value_type};\n"
                    config_class_code += "}\n"
                except Exception as e:
                    pass

                ts_type = f"Array<{innerClassName}>"
            field_lists.append({
                "name": headerRow[key],
                "type": ts_type
            })

        config_name = str.upper(class_name[0])+ class_name[1:]
        config_class_code += f'''
        export default class {config_name}Data extends BaseConfigItem {{
          public static fileName:string = "{class_name}Data";
        '''
        for field in field_lists:
            config_class_code += f'''
            private _{field["name"]}: {field["type"]};
            public get {field["name"]}():{field["type"]} {{return this._{field["name"]};}}
            '''
        config_class_code += f'''
        public constructor(data:any) {{
            super(data);
        '''
        for field in field_lists:
            config_class_code += f"        this._{field['name']} = data['{field['name']}']\n"
        config_class_code += "    }\n}\n"

        with open(target_code_path, 'w', encoding='utf-8') as f:
            f.write(config_class_code)
            print(f"Code file has been successfully generated for file {class_name}")

    @staticmethod
    def exportData(source_directory):
        # 定义生成代码文件的目标目录，按需修改
        excel_files = glob.glob(os.path.join(source_directory, "*.xlsx"))
        try:
            for excel_file in excel_files:
                json_data = ExportUtil.excel2json(excel_file)
                # ensure_ascii表示禁止中文转义
                json_str = json.dumps(json_data, indent=4, ensure_ascii=False)

                file_name_without_ext = os.path.splitext(os.path.basename(excel_file))[0]
                json_file_path = os.path.join(target_json_directory, file_name_without_ext + "Data.json")

                with open(json_file_path, 'w', encoding='utf-8') as f:
                    f.write(json_str)
                print(f"JSON data has been successfully written to {json_file_path} for file {excel_file}")
        except Exception as e:
            print(e)

if __name__ == '__main__':
    ExportUtil.exportData(source_directory)