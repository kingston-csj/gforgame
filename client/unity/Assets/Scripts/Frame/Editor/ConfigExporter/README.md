### excel配置说明

- HEADER行以上的记录程序不会读写，策划可自行添加，例如增加注释，字段类型等等
- HEADER行标记每个字段的名称，程序从HEADER所在行的下一行开始读取
- 程序读取到END所在的行结束，END行下面，即使有数据，程序也不会读取
- 如果没有End标记，程序会一直读取到最后一行记录，如果文件出现空白行，记录可能会有异常！！

  ![Image](config_export "配置格式")

- EXPORT所在的行为可选项,没有则代表所有字段都导出
- SERVER表示该字段为服务器使用，客户端不需要
- CLIENT表示该字段为客户端使用，服务器不需要
- BOTH表示该字段为服务器和客户端都需要
- 空白表示服务器和客户端均不需要,仅作策划备注

### excel读取工具
使用了第三方库ExcelDataReader，将下载的dll放在ThirdParty目录下
https://github.com/ExcelDataReader/ExcelDataReader#usage