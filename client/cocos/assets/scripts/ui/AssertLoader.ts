import { Prefab, resources } from 'cc';

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
      console.log('预制体加载成功:', prefab);
      callback(null, prefab);
    });
  }
}
