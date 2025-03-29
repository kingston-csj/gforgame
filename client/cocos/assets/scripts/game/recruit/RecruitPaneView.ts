import { _decorator, Node } from 'cc';

import { BaseUiView } from '../../ui/BaseUiView';

const { ccclass, property } = _decorator;

@ccclass('RecruitPaneView')
export class RecruitPaneView extends BaseUiView {
  @property(Node)
  oneBtn: Node;

  @property(Node)
  tenBtn: Node;

  @property(Node)
  closeBtn: Node;

  protected start(): void {
    this.registerClickEvent(this.oneBtn, () => this.onRecruitBtnClick(1), this);
    this.registerClickEvent(this.tenBtn, () => this.onRecruitBtnClick(10), this);
    this.registerClickEvent(this.closeBtn, () => this.onCloseBtnClick(), this);
  }

  onRecruitBtnClick(times: number) {
    this.node.emit('recruitBtnClick', times);
  }

  onCloseBtnClick() {
    this.node.emit('closeBtnClick');
    this.hide();
  }
}
