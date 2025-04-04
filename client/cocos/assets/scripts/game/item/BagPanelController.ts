import { _decorator } from 'cc';
import { BaseController } from '../../ui/BaseController';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
import { RecruitModel } from '../recruit/RecruitModel';
import { BagPanelView } from './BagPanelView';

const { ccclass, property } = _decorator;

@ccclass('BagPanelController')
export class BagPanelController extends BaseController {
  private static instance: BagPanelController;

  private static creatingPromise: Promise<BagPanelController> | null = null;

  @property(BagPanelView)
  recruitPaneView: BagPanelView | null = null;

  recruitModel: RecruitModel = RecruitModel.instance;

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

  private static getInstance(): Promise<BagPanelController> {
    if (this.instance) {
      return Promise.resolve(this.instance);
    }
    if (this.creatingPromise) {
      return this.creatingPromise;
    }
    this.creatingPromise = new Promise((resolve) => {
      UiViewFactory.createUi(R.Prefabs.BagPane, LayerIdx.layer4, (ui: BagPanelController) => {
        this.instance = ui;
        this.creatingPromise = null;
        resolve(ui);
      });
    });
    return this.creatingPromise;
  }
}
