import { _decorator, Component, instantiate, Node } from 'cc';
import ResourceItem from './ResourceItem';
import { LayerIdx } from './LayerIds';
import AssetLoader from './AssertLoader';
import UiContext from './UiContext';
import { BaseController } from './BaseController';

const { ccclass, property } = _decorator;

@ccclass('UiViewFactory')
export default class UiViewFactory extends Component {
  public static createUi(ui: ResourceItem, layer: LayerIdx, callback: Function) {
    // 使用AssetLoader加载预制体
    AssetLoader.loadPrefab(ui.path, (err, prefab) => {
      if (err) {
        console.error('加载UI预制体失败:', err);
        return;
      }

      // 实例化预制体
      const node = instantiate(prefab);
      if (!node) {
        console.error('实例化UI预制体失败');
        callback(new Error('实例化UI预制体失败'), null);
        return;
      }

      const root = UiContext.getLayer(layer);
      root.addChild(node);
      const uiView = node.getComponent(BaseController);
      if (uiView) {
        uiView.scheduleOnce(() => {
          if (callback) {
            callback(uiView);
          }
        }, 0.2);
      }
    });
  }
}
