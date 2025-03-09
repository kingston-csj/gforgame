import { _decorator, Component, Node, EditBox, Button, director } from 'cc';
import GameContext from '../GameContext';

import { LayerIdx } from '../ui/LayerIds';

import R from '../ui/R';
import UiContext from '../ui/UiContext';
import UiView from '../ui/UiView';
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
    UiContext.init(this.layer1, this.layer2, this.layer3, this.layer4, this.layer5);
    GameContext.instance.WebSocketClient.connect('ws://127.0.0.1:9527/ws');

    UiView.createUi(R.loginPane, LayerIdx.layer1, () => {});
  }
}
