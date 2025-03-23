import { Prefab, resources, SpriteFrame, ImageAsset, assetManager } from 'cc';

export default class AssetLoader {
  /**
   * 根据路径加载预制体
   * @param prefabPath 预制体在 resources 文件夹下的相对路径
   * @param callback 加载完成后的回调函数，参数为加载结果（错误信息或预制体对象）
   */
  static loadPrefab(
    prefabPath: string,
    callback: (err: Error | null, prefab: Prefab | null) => void
  ) {
    resources.load(prefabPath, Prefab, (err, prefab: Prefab) => {
      if (err) {
        console.error('加载预制体失败:', err);
        callback(err, null);
        return;
      }
      callback(null, prefab);
    });
  }

  /**
   * 根据路径加载精灵帧
   * @param spritePath 精灵帧在 resources 文件夹下的相对路径
   * @param callback 加载完成后的回调函数，参数为加载结果（错误信息或精灵帧对象）
   */
  static loadSpriteFrame(
    spritePath: string,
    callback: (err: Error | null, spriteFrame: SpriteFrame | null) => void
  ) {
    resources.load(spritePath, SpriteFrame, (err, spriteFrame: SpriteFrame) => {
      if (err) {
        console.error('加载精灵帧失败:', err);
        callback(err, null);
        return;
      }
      callback(null, spriteFrame);
    });
  }

  /**
   * 根据路径加载图片资源
   * @param imagePath 图片在 resources 文件夹下的相对路径
   * @param callback 加载完成后的回调函数，参数为加载结果（错误信息或图片资源对象）
   */
  static loadImageAsset(
    imagePath: string,
    callback: (err: Error | null, imageAsset: ImageAsset | null) => void
  ) {
    resources.load(imagePath, ImageAsset, (err, imageAsset: ImageAsset) => {
      if (err) {
        console.error('加载图片资源失败:', err);
        callback(err, null);
        return;
      }
      callback(null, imageAsset);
    });
  }
}
