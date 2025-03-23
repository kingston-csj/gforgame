import { _decorator, Component, Node, EditBox, Button, director } from 'cc';
import GameContext from '../../GameContext';

import { UIViewController } from '../../ui/UiViewController';

import UiView from '../../ui/UiView';
import { LayerIdx } from '../../ui/LayerIds';
import { ReqHeroRecruit } from '../../net/ReqHeroRecruit';
import { ResHeroRecruit } from '../../net/ResHeroRecruit';
import R from '../../ui/R';
import { RecruitSettleModel } from '../../model/RecruitSettleModel';
import { RecruitSettlePaneController } from './RecruitSettlePaneController';
const { ccclass, property } = _decorator;

@ccclass('RecruitPaneController')
export class RecruitPaneController extends UIViewController {
  @property(Node)
  oneBtn: Node;

  @property(Node)
  tenBtn: Node;

  @property(Node)
  closeBtn: Node;

  private static instance: RecruitPaneController;

  protected start(): void {
    this.registerClickEvent(this.oneBtn, () => this.onRecruitBtnClick(1), this);
    this.registerClickEvent(this.tenBtn, () => this.onRecruitBtnClick(10), this);
    this.registerClickEvent(this.closeBtn, () => this.onCloseBtnClick(), this);
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
    console.log('关闭');
    this.hide();
  }

  public static openUi() {
    if (RecruitPaneController.instance) {
      RecruitPaneController.instance.display();
    } else {
      RecruitPaneController.instance = new RecruitPaneController();

      UiView.createUi(R.recruitPane, LayerIdx.layer4, (ui: RecruitPaneController) => {
        RecruitPaneController.instance = ui;
        ui.display();
      });
    }
  }
}
