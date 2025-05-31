import { RedDotComponent } from './RedDotCompoent';
import { RedDotNode } from './RedDotNode';

export class RedDotManager {
  private static _instance: RedDotManager = new RedDotManager();

  // 根节点
  private _root: RedDotNode = null;

  public static get instance(): RedDotManager {
    return RedDotManager._instance;
  }

  constructor() {
    this._root = new RedDotNode();
  }

  public binding(path: string, ui: RedDotComponent) {
    const node = this.getOrCreateNode(path);
    node.ui = ui;
    ui.updateScore(node.score);
  }

  private getOrCreateNode(path: string): RedDotNode {
    const paths = path.split('/');
    let current = this._root;
    for (const p of paths) {
      current = current.addChild(p);
    }
    return current;
  }

  public updateScore(path: string, score: number) {
    const node = this.getOrCreateNode(path);
    node.score = score;
    // 向上回溯，更新父节点分数
    let parent = node.parent;
    while (parent) {
      let parentScore = 0;
      for (const child of parent.children.values()) {
        parentScore += child.score;
      }
      parent.score = parentScore;
      parent = parent.parent;
    }

    // 从当前节点开始，向上更新UI
    let current = node;
    while (current) {
      current.ui?.updateScore(current.score);
      current = current.parent;
    }
  }
}
