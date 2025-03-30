import { _decorator } from 'cc';
import { ConfigI18nContainer } from '../../data/config/container/ConfigI18nContainer';
import { BaseController } from '../../ui/BaseController';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
import { TipsView } from './TipsView';
const { ccclass, property } = _decorator;

@ccclass('TipsPaneController')
export class TipsPaneController extends BaseController {
  @property(TipsView)
  tipsView: TipsView | null = null;

  private static instance: TipsPaneController;

  private static creatingPromise: Promise<TipsPaneController> | null = null;

  start() {
    this.initView(this.tipsView);
  }

  public static openUi(code: number) {
    this.getInstance().then((controller) => {
      if (controller.tipsView) {
        let tips = ConfigI18nContainer.getInstance().getRecord(code).content;
        controller.tipsView.setTips(tips);
        controller.tipsView.display();
        // 显示后多等X秒
        setTimeout(() => {
          controller.tipsView.hide();
        }, 2000);
      }
    });
  }

  private static getInstance(): Promise<TipsPaneController> {
    if (this.instance) {
      return Promise.resolve(this.instance);
    }
    if (this.creatingPromise) {
      return this.creatingPromise;
    }
    this.creatingPromise = new Promise((resolve) => {
      UiViewFactory.createUi(R.Prefabs.TipsPane, LayerIdx.layer5, (ui: TipsPaneController) => {
        this.instance = ui;
        this.creatingPromise = null;
        resolve(ui);
      });
    });
    return this.creatingPromise;
  }
}
