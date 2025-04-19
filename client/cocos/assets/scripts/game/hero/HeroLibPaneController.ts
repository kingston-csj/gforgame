import { _decorator } from 'cc';
import { BaseController } from '../../ui/BaseController';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
import { HeroLibView } from './HeroLibView';

const { ccclass, property } = _decorator;

@ccclass('HeroLibPaneController')
export class HeroLibPaneController extends BaseController {
  private static instance: HeroLibPaneController;

  private static creatingPromise: Promise<HeroLibPaneController> | null = null;

  @property(HeroLibView)
  private mainView: HeroLibView = null;

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

  private static getInstance(): Promise<HeroLibPaneController> {
    if (this.instance) {
      return Promise.resolve(this.instance);
    }
    if (this.creatingPromise) {
      return this.creatingPromise;
    }
    this.creatingPromise = new Promise((resolve) => {
      UiViewFactory.createUi(
        R.Prefabs.HeroLibPane,
        LayerIdx.layer4,
        (ui: HeroLibPaneController) => {
          this.instance = ui;
          this.creatingPromise = null;
          resolve(ui);
        }
      );
    });
    return this.creatingPromise;
  }
}
