import { _decorator, Component, Node } from 'cc';
import GameContext from '../GameContext';

import { LayerIdx } from '../ui/LayerIds';

import ConfigHeroContainer from '../data/config/container/ConfigHeroContainer';
import ConfigItemContainer from '../data/config/container/ConfigItemContainer';
import { MessageDispatch } from '../MessageDispatch';
import R from '../ui/R';
import UiContext from '../ui/UiContext';
import UiViewFactory from '../ui/UiViewFactory';

const { ccclass, property } = _decorator;

@ccclass('Game')
export class LoginScene extends Component {
  @property(Node)
  layer1: Node;

  @property(Node)
  layer2: Node;

  @property(Node)
  layer3: Node;

  @property(Node)
  layer4: Node;

  @property(Node)
  layer5: Node;

  start() {
    // 初始化资源工厂
    UiViewFactory.createUi(R.Prefabs.AssetResourceFactory, LayerIdx.layer1, () => {});

    // 加载所有的配置数据
    ConfigItemContainer.getInstance();
    ConfigHeroContainer.getInstance();

    // 挂载备份节点
    UiContext.init(this.layer1, this.layer2, this.layer3, this.layer4, this.layer5);

    // 初始化消息监听
    MessageDispatch.init();

    // 连接服务器
    GameContext.instance.WebSocketClient.connect('ws://127.0.0.1:9527/ws');

    // 创建登录界面
    UiViewFactory.createUi(R.Prefabs.LoginPane, LayerIdx.layer1, () => {});
  }
}
