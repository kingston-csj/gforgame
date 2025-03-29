import { _decorator, Component, JsonAsset, SpriteAtlas } from 'cc';
const { ccclass, property } = _decorator;
@ccclass('AssetResourceFactory')
export default class AssetResourceFactory extends Component {
  private static _instance: AssetResourceFactory;

  @property({
    type: [SpriteAtlas],
  })
  public spriteAtlas: Array<SpriteAtlas> = [];

  @property({
    type: [JsonAsset],
  })
  public configs: Array<JsonAsset> = [];

  public static get instance() {
    if (!this._instance) {
      this._instance = new AssetResourceFactory();
    }
    return this._instance;
  }

  protected onLoad(): void {
    AssetResourceFactory._instance = this;
  }

  public getConfig(name: string): JsonAsset {
    return this.configs.find((config) => config.name === name);
  }

  public getSpriteAtlas(name: string): SpriteAtlas {
    return this.spriteAtlas.find((spriteAtlas) => spriteAtlas.name === name);
  }
}
