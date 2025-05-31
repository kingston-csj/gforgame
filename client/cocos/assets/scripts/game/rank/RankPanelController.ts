import { _decorator } from 'cc';
import { BaseController } from '../../frame/mvc/BaseController';

import R from '../../ui/R';

import { LayerIdx } from '../../ui/LayerIds';
import UiViewFactory from '../../ui/UiViewFactory';
import { RankView } from './RankView';
const { ccclass, property } = _decorator;

@ccclass('RankPanelController')
export class RankPanelController extends BaseController {
  private static instance: RankPanelController;
  private static creatingPromise: Promise<RankPanelController> | null = null;

  @property(RankView)
  public view: RankView;

  public static openUi() {
    this.getInstance().then((controller) => {
      controller.view.display();
    });
  }

  private static getInstance(): Promise<RankPanelController> {
    if (this.instance) {
      return Promise.resolve(this.instance);
    }

    if (this.creatingPromise) {
      return this.creatingPromise;
    }

    this.creatingPromise = new Promise((resolve) => {
      UiViewFactory.createUi(R.Prefabs.RankPane, LayerIdx.layer4, (ui: RankPanelController) => {
        this.instance = ui;
        this.creatingPromise = null;
        resolve(ui);
      });
    });
    return this.creatingPromise;
  }
}
