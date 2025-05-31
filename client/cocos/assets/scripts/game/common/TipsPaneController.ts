import { _decorator } from 'cc';
import { ConfigContext } from '../../data/config/container/ConfigContext';
import { BaseController } from '../../frame/mvc/BaseController';
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

  public static showI18nContent(code: number) {
    this.getInstance().then((controller) => {
      if (controller.tipsView) {
        let tips = ConfigContext.configI18nContainer.getRecord(code).content;
        controller.tipsView.setTips(tips);
        controller.tipsView.display();
        // 显示后多等X秒
        setTimeout(() => {
          controller.tipsView.hide();
        }, 2000);
      }
    });
  }

  public static showStringContent(content: string) {
    this.getInstance().then((controller) => {
      if (controller.tipsView) {
        controller.tipsView.setTips(content);
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
