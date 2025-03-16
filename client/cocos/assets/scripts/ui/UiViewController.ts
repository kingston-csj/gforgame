import { _decorator, Button, Component, instantiate, Node, Prefab, Toggle } from 'cc';
const { ccclass, property } = _decorator;
@ccclass('UI_View')
export class UIViewController extends Component {
  public display() {
    this.node.active = true;
    this.onDisplay();
  }

  protected onDisplay() {}

  public hide() {
    this.node.active = false;
    this.onHide();
  }

  protected onHide() {}
}
