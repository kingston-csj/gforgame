import { _decorator, Component } from 'cc';

import UiViewFactory from '../../ui/UiViewFactory';
import { LayerIdx } from '../../ui/LayerIds';

import { ResHeroRecruit } from '../../net/ResHeroRecruit';
import R from '../../ui/R';
import { RecruitSettleModel } from './RecruitSettleModel';

import { RecruitSettleView } from './RecruitSettleView';
import { BaseController } from '../../ui/BaseController';
const { ccclass, property } = _decorator;

@ccclass('RecruitSettlePaneController')
export class RecruitSettlePaneController extends BaseController {
  @property(RecruitSettleView)
  public recruitSettleView: RecruitSettleView | null = null;

  private recruitSettleModel: RecruitSettleModel = RecruitSettleModel.getInstance();

  private static creatingPromise: Promise<RecruitSettlePaneController> | null = null;

  private static instance: RecruitSettlePaneController;

  protected start(): void {
    this.recruitSettleView.node.on('recruitBtnClick', this.onRecruitBtnClick, this);
    this.recruitSettleView.node.on('closeBtnClick', this.onCloseBtnClick, this);
  }

  onRecruitBtnClick(times: number) {
    this.recruitSettleModel.doRecruit(times).then((msg: ResHeroRecruit) => {
      this.recruitSettleModel.setRewardItems(msg.rewardInfos);
      this.recruitSettleView.display();
    });
  }

  onCloseBtnClick() {
    this.recruitSettleView.hide();
  }

  public static openUi() {
    this.getInstance().then((controller) => {
      if (controller.recruitSettleView) {
        controller.recruitSettleView.display();
      }
    });
  }

  private static getInstance(): Promise<RecruitSettlePaneController> {
    if (this.instance) {
      return Promise.resolve(this.instance);
    }
    if (this.creatingPromise) {
      return this.creatingPromise;
    }

    this.creatingPromise = new Promise((resolve) => {
      UiViewFactory.createUi(
        R.recruitSettlePane,
        LayerIdx.layer5,
        (ui: RecruitSettlePaneController) => {
          this.instance = ui;
          this.creatingPromise = null;
          resolve(ui);
        }
      );
    });
    return this.creatingPromise;
  }
}
