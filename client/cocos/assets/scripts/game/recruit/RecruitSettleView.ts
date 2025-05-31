import { _decorator, instantiate, Node, Prefab } from 'cc';

import { BaseUiView } from '../../frame/mvc/BaseUiView';

import { RewardItem } from '../reward/RewardItem';
import { RecruitSettleModel } from './RecruitSettleModel';

const { ccclass, property } = _decorator;

@ccclass('RecruitSettleView')
export class RecruitSettleView extends BaseUiView {
  @property(Node)
  public itemContainer: Node;

  @property(Prefab)
  public itemPrefab: Prefab;

  @property(Node)
  public againBtn: Node;

  @property(Node)
  public okBtn: Node;

  protected start(): void {
    this.registerClickEvent(this.againBtn, () => this.onRecruitBtnClick(1), this);
    this.registerClickEvent(this.okBtn, this.hide, this);
  }

  onRecruitBtnClick(times: number) {
    this.node.emit('recruitBtnClick', times);
  }

  protected onDisplay() {
    this.showItems();
  }

  private showItems() {
    this.itemContainer.removeAllChildren();
    let rewardItems = RecruitSettleModel.getInstance().getRewardItems();
    for (let i = 0; i < rewardItems.length; i++) {
      let itemUi = instantiate(this.itemPrefab);
      itemUi.setParent(this.itemContainer);
      itemUi.getComponent(RewardItem).fillData(rewardItems[i]);
    }
  }
}
