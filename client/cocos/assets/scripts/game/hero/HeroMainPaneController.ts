import { _decorator } from 'cc';
import { BaseController } from '../../ui/BaseController';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
import { HeroMainView } from './HeroMainView';

const { ccclass, property } = _decorator;

@ccclass('HeroMainPaneController')
export class HeroMainPaneController extends BaseController {
  private static instance: HeroMainPaneController;

  private static creatingPromise: Promise<HeroMainPaneController> | null = null;

  @property(HeroMainView)
  private mainView: HeroMainView | null = null;

  private constructor() {
    super();
  }

  public static openUi() {
    this.getInstance().then((controller) => {
      if (controller.mainView) {
        controller.mainView.display();
      }
    });
  }

  public static closeUi() {
    if (!this.instance) {
      return Promise.resolve();
    }
    return this.getInstance().then((controller) => {
      if (controller.mainView) {
        controller.mainView.hide();
      }
    });
  }

  protected start(): void {
    this.initView(this.mainView);
  }

  private static getInstance(): Promise<HeroMainPaneController> {
    if (this.instance) {
      return Promise.resolve(this.instance);
    }
    if (this.creatingPromise) {
      return this.creatingPromise;
    }
    this.creatingPromise = new Promise((resolve) => {
      UiViewFactory.createUi(
        R.Prefabs.HeroMainPane,
        LayerIdx.layer2,
        (ui: HeroMainPaneController) => {
          this.instance = ui;
          this.creatingPromise = null;
          resolve(ui);
        }
      );
    });
    return this.creatingPromise;
  }
}
