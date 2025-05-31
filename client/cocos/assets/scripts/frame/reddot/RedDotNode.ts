import { _decorator } from 'cc';
import { RedDotComponent } from './RedDotCompoent';
const { ccclass, property } = _decorator;

@ccclass('RedDotNode')
export class RedDotNode {
  /**
   * 分数
   */
  private _score: number = 0;

  private _parent: RedDotNode;

  private _children: Map<string, RedDotNode> = new Map<string, RedDotNode>();

  private _ui: RedDotComponent = null;

  public addChild(name: string): RedDotNode {
    if (this._children.has(name)) {
      return this._children.get(name);
    }
    const node = new RedDotNode();
    this._children.set(name, node);
    node._parent = this;
    return node;
  }

  public set ui(ui: RedDotComponent) {
    this._ui = ui;
  }

  public get ui(): RedDotComponent {
    return this._ui;
  }

  public set score(score: number) {
    this._score = score;
  }

  public get score(): number {
    return this._score;
  }

  public get parent(): RedDotNode {
    return this._parent;
  }

  public get children(): Map<string, RedDotNode> {
    return this._children;
  }
}
