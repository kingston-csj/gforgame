import { _decorator, Button, Component, Node, Toggle } from 'cc';
const { ccclass, property } = _decorator;
@ccclass('BaseUiView')
export class BaseUiView extends Component {
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

  public isShow(): boolean {
    return this.node.active;
  }

  /**
   * 为节点注册点击事件
   * @param node 需要注册点击事件的节点
   * @param callback 点击回调函数
   * @param target 回调函数的this指向，默认为当前组件
   */
  protected registerClickEvent(node: Node, callback: () => void, target: any = this) {
    // 检查节点上是否有Button组件
    const button = node.getComponent(Button);
    if (button) {
      // 如果有Button组件，注册Button的点击事件
      button.node.on(Button.EventType.CLICK, callback, target);
      return;
    }

    // 检查节点上是否有Toggle组件
    const toggle = node.getComponent(Toggle);
    if (toggle) {
      // 如果有Toggle组件，注册Toggle的点击事件
      toggle.node.on(Toggle.EventType.CLICK, callback, target);
      return;
    }

    // 如果既没有Button也没有Toggle，作为普通节点注册点击事件
    node.on(Node.EventType.TOUCH_END, callback, target);
  }
}
