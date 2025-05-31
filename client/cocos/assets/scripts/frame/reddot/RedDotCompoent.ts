import { _decorator, Component, Label, Node } from 'cc';
import { RedDotManager } from './RedDotManager';
const { ccclass, property } = _decorator;

@ccclass('RedDotComponent')
export class RedDotComponent extends Component {
  @property
  public path: string = '';

  @property(Label)
  public scoreLabel: Label = null;

  @property(Node)
  public circle: Node = null;

  /**
   * 是否显示数字
   */
  @property
  public showNumber: boolean = false;

  public updateScore(score: number) {
    this.circle.active = score > 0;
    if (this.showNumber && this.scoreLabel) {
      this.scoreLabel.node.active = true;
      this.scoreLabel.string = score.toString();
    }
  }

  protected onLoad() {
    if (this.path) {
      RedDotManager.instance.binding(this.path, this);
    }
  }
}
