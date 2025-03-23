import { _decorator, Component, Node, EditBox, Button, director, instantiate, Prefab } from 'cc';
import GameContext from '../../GameContext';

import { UIViewController } from '../../ui/UiViewController';

import UiView from '../../ui/UiView';
import { LayerIdx } from '../../ui/LayerIds';
import { ReqHeroRecruit } from '../../net/ReqHeroRecruit';
import { ResHeroRecruit } from '../../net/ResHeroRecruit';
import R from '../../ui/R';
import { RecruitSettleModel } from '../../model/RecruitSettleModel';
import { RewardItem } from '../reward/RewardItem';
const { ccclass, property } = _decorator;

@ccclass('RecruitSettlePaneController')
export class RecruitSettlePaneController extends UIViewController {
  @property(Node)
  public itemContainer: Node;

  @property(Prefab)
  public itemPrefab: Prefab;

  @property(Node)
  public againBtn: Node;

  @property(Node)
  public okBtn: Node;

  private static instance: RecruitSettlePaneController;

  protected start(): void {
    this.registerClickEvent(this.againBtn, () => this.onRecruitBtnClick(1), this);
    this.registerClickEvent(this.okBtn, this.hide, this);
  }

  onRecruitBtnClick(time: number) {
    GameContext.instance.WebSocketClient.sendMessage(
      ReqHeroRecruit.cmd,
      {
        times: time,
      },
      (msg: ResHeroRecruit) => {
        RecruitSettleModel.getInstance().setRewardItems(msg.rewardInfos);
        RecruitSettlePaneController.openUi();
      }
    );
  }

  onCloseBtnClick() {
    this.hide();
  }

  public static openUi() {
    if (RecruitSettlePaneController.instance) {
      RecruitSettlePaneController.instance.display();
    } else {
      RecruitSettlePaneController.instance = new RecruitSettlePaneController();

      UiView.createUi(R.recruitSettlePane, LayerIdx.layer5, (ui: RecruitSettlePaneController) => {
        RecruitSettlePaneController.instance = ui;
        ui.display();
      });
    }
  }

  protected onDisplay() {
    this.showItems();
  }

  private showItems() {
    let rewardItems = RecruitSettleModel.getInstance().getRewardItems();
    for (let i = 0; i < rewardItems.length; i++) {
      let itemUi = instantiate(this.itemPrefab);
      itemUi.setParent(this.itemContainer);

      itemUi.getComponent(RewardItem).fillData(rewardItems[i]);
    }
  }
}
