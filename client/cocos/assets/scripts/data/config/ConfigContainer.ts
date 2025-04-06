import AssetResourceFactory from '../../ui/AssetResourceFactory';

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
    let jsonAsset = AssetResourceFactory.instance.getConfig(this._fileName);
    let jsonData = jsonAsset.json;
    if (!Array.isArray(jsonData)) {
      console.error(`配置文件 ${this._fileName} 格式错误，应该是数组`);
      return;
    }
    for (const item of jsonData) {
      try {
        const config = new this._meta(item);
        this._datas.set(config.id, config as T);
      } catch (e) {
        console.error(`解析配置项失败:`, item, e);
      }
    }

    console.log(`配置文件 ${this._fileName} 加载完成，共 ${this._datas.size} 条数据`);
  }

  /**
   * 加载配置文件后钩子，子类可以重写此方法，存储一些二级缓存
   */
  public afterLoad() {}

  // 获取单个配置
  public getRecord(id: number): T | undefined {
    return this._datas.get(id);
  }

  // 获取所有配置
  public getAllRecords(): T[] {
    return Array.from(this._datas.values());
  }
}
