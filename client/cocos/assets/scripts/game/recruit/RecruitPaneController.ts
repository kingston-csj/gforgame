import { _decorator, Component } from 'cc';

import { ResHeroRecruit } from '../../net/ResHeroRecruit';

import { RecruitSettleModel } from './RecruitSettleModel';
import { RecruitSettlePaneController } from './RecruitSettlePaneController';
import { RecruitPaneView } from './RecruitPaneView';
import { RecruitModel } from './RecruitModel';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
import { BaseController } from '../../ui/BaseController';
const { ccclass, property } = _decorator;

@ccclass('RecruitPaneController')
export class RecruitPaneController extends BaseController {
  private static instance: RecruitPaneController;

  private static creatingPromise: Promise<RecruitPaneController> | null = null;

  @property(RecruitPaneView)
  recruitPaneView: RecruitPaneView | null = null;

  recruitModel: RecruitModel = RecruitModel.instance;

  start() {
    this.initView(this.recruitPaneView);
  }

  protected bindViewEvents() {
    this.recruitPaneView.node.on('recruitBtnClick', this.onRecruitBtnClick, this);
    this.recruitPaneView.node.on('closeBtnClick', this.onCloseBtnClick, this);
  }

  onRecruitBtnClick(times: number) {
    this.recruitModel.doRecruit(times).then((msg: ResHeroRecruit) => {
      RecruitSettleModel.getInstance().setRewardItems(msg.rewardInfos);
      RecruitSettlePaneController.openUi();
    });
  }

  onCloseBtnClick() {
    this.recruitPaneView.hide();
  }

  public static openUi() {
    this.getInstance().then((controller) => {
      if (controller.recruitPaneView) {
        controller.recruitPaneView.display();
      }
    });
  }

  private static getInstance(): Promise<RecruitPaneController> {
    if (this.instance) {
      return Promise.resolve(this.instance);
    }
    if (this.creatingPromise) {
      return this.creatingPromise;
    }
    this.creatingPromise = new Promise((resolve) => {
      UiViewFactory.createUi(R.recruitPane, LayerIdx.layer4, (ui: RecruitPaneController) => {
        this.instance = ui;
        this.creatingPromise = null;
        resolve(ui);
      });
    });
    return this.creatingPromise;
  }
}
