import { _decorator } from 'cc';

import { ResHeroRecruit } from '../../net/protocol/ResHeroRecruit';

import { BaseController } from '../../frame/mvc/BaseController';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
import { TipsPaneController } from '../common/TipsPaneController';

import GameConstants from '../constants/GameConstants';
import BagpackModel from '../item/BagpackModel';
import { RecruitModel } from './RecruitModel';
import { RecruitPaneView } from './RecruitPaneView';
import { RecruitSettleModel } from './RecruitSettleModel';
import { RecruitSettlePaneController } from './RecruitSettlePaneController';
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
    let ownItem = BagpackModel.getInstance().getItemByModelId(GameConstants.Item.RECRUIT_ID);
    // if (!ownItem || ownItem.count < times) {
    //   return;
    // }
    this.recruitModel.doRecruit(times).then((msg: ResHeroRecruit) => {
      if (msg.code === 0) {
        RecruitSettleModel.getInstance().setRewardItems(msg.rewardInfos);
        RecruitSettlePaneController.openUi();
      } else {
        TipsPaneController.showI18nContent(msg.code);
      }
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

  public static closeUi() {
    if (!this.instance) {
      return Promise.resolve();
    }
    return this.getInstance().then((controller) => {
      if (controller.recruitPaneView) {
        controller.recruitPaneView.hide();
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
      UiViewFactory.createUi(
        R.Prefabs.RecruitPane,
        LayerIdx.layer4,
        (ui: RecruitPaneController) => {
          this.instance = ui;
          this.creatingPromise = null;
          resolve(ui);
        }
      );
    });
    return this.creatingPromise;
  }
}
