import { _decorator } from 'cc';

import { LayerIdx } from '../../ui/LayerIds';
import UiViewFactory from '../../ui/UiViewFactory';

import { ResHeroRecruit } from '../../net/protocol/ResHeroRecruit';
import R from '../../ui/R';
import { RecruitSettleModel } from './RecruitSettleModel';

import { BaseController } from '../../ui/BaseController';
import { TipsPaneController } from '../common/TipsPaneController';
import { RecruitSettleView } from './RecruitSettleView';
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
      if (msg.code === 0) {
        RecruitSettleModel.getInstance().setRewardItems(msg.rewardInfos);
        RecruitSettlePaneController.openUi();
      } else {
        TipsPaneController.openUi(msg.code);
      }
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
        R.Prefabs.RecruitSettlePane,
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
