import { _decorator } from 'cc';
import { BaseController } from '../../ui/BaseController';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
import { GmPaneView } from './GmPaneView';

const { ccclass, property } = _decorator;

@ccclass('GmPaneController')
export class GmPaneController extends BaseController {
  private static instance: GmPaneController;

  @property(GmPaneView)
  private gmPaneView: GmPaneView = null!;

  private static creatingPromise: Promise<GmPaneController> | null = null;

  protected start(): void {
    this.initView(this.gmPaneView);
  }

  private static getInstance(): Promise<GmPaneController> {
    if (this.instance) {
      return Promise.resolve(this.instance);
    }
    if (this.creatingPromise) {
      return this.creatingPromise;
    }
    this.creatingPromise = new Promise((resolve) => {
      UiViewFactory.createUi(R.Prefabs.GmPane, LayerIdx.layer2, (ui: GmPaneController) => {
        this.instance = ui;
        this.creatingPromise = null;
        resolve(ui);
      });
    });
    return this.creatingPromise;
  }

  public static switchDisplay() {
    GmPaneController.getInstance().then((controller) => {
      if (controller.gmPaneView.isShow()) {
        controller.gmPaneView.hide();
      } else {
        controller.gmPaneView.display();
      }
    });
  }
}
