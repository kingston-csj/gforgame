import { JsonAsset, resources } from 'cc';

export default class ConfigContainer<T> {
  protected _datas: Map<number, T> = new Map<number, T>();
  private _meta: any;
  private _fileName: string;

  public constructor(meta: any, fileName: string) {
    this._meta = meta;
    this._fileName = fileName;
    this.loadConfig();
  }

  private loadConfig() {
    // 加载JSON配置文件，路径相对于resources目录
    resources.load(`config/${this._fileName}`, JsonAsset, (err, jsonAsset) => {
      if (err) {
        console.error(`加载配置文件 ${this._fileName} 失败:`, err);
        return;
      }

      const jsonData = jsonAsset.json;
      if (!Array.isArray(jsonData)) {
        console.error(`配置文件 ${this._fileName} 格式错误，应该是数组`);
        return;
      }

      // 解析每一条配置数据
      for (const item of jsonData) {
        try {
          // 使用传入的类型创建实例
          const config = new this._meta(item);
          this._datas.set(config.id, config as T);
          // 使用id作为key存储到Map中
          if ('id' in config) {
            this._datas.set(config.id, config as T);
          } else {
            console.error(`配置项缺少id字段:`, item);
          }
        } catch (e) {
          console.error(`解析配置项失败:`, item, e);
        }
      }

      console.log(`配置文件 ${this._fileName} 加载完成，共 ${this._datas.size} 条数据`);
      this.afterLoad();
    });
  }

  /**
   * 加载配置文件后钩子，子类可以重写此方法，存储一些二级缓存
   */
  protected afterLoad() {}

  // 获取单个配置
  public getRecord(id: number): T | undefined {
    return this._datas.get(id);
  }

  // 获取所有配置
  public getAllRecords(): T[] {
    return Array.from(this._datas.values());
  }

  // 静态创建方法
  public static create<T>(meta: any, fileName: string): ConfigContainer<T> {
    return new ConfigContainer<T>(meta, fileName);
  }
}
