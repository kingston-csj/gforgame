import { _decorator } from 'cc';
import { BaseController } from '../../frame/mvc/BaseController';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
import { BuZhenView } from './BuZhenView';
const { ccclass, property } = _decorator;

@ccclass('BuZhenPaneController')
export class BuZhenPaneController extends BaseController {
  private static instance: BuZhenPaneController;

  private static creatingPromise: Promise<BuZhenPaneController> | null = null;

  @property(BuZhenView)
  public view: BuZhenView;

  public static openUi() {
    this.getInstance().then((controller) => {
      controller.view.display();
    });
  }

  protected start(): void {
    this.initView(this.view);
  }

  private static getInstance(): Promise<BuZhenPaneController> {
    if (this.instance) {
      return Promise.resolve(this.instance);
    }

    if (this.creatingPromise) {
      return this.creatingPromise;
    }
    this.creatingPromise = new Promise((resolve) => {
      UiViewFactory.createUi(R.Prefabs.BuZhenPane, LayerIdx.layer4, (ui: BuZhenPaneController) => {
        this.instance = ui;
        this.creatingPromise = null;
        resolve(ui);
      });
    });
    return this.creatingPromise;
  }

  public static refreshLineupHeros(): void {
    this.getInstance().then((controller) => {
      controller.view.showLineupHeros();
    });
  }
}
