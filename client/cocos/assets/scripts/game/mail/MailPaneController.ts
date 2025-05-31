import { _decorator } from 'cc';
import { BaseController } from '../../frame/mvc/BaseController';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
import { MailView } from './MailView';

const { ccclass, property } = _decorator;

@ccclass('MailPaneController')
export class MailPaneController extends BaseController {
  private static instance: MailPaneController;
  private static creatingPromise: Promise<MailPaneController> | null = null;

  @property(MailView)
  public view: MailView;

  public static openUi() {
    this.getInstance().then((controller) => {
      controller.view.display();
    });
  }

  protected start(): void {
    this.initView(this.view);
  }

  private static getInstance(): Promise<MailPaneController> {
    if (this.instance) {
      return Promise.resolve(this.instance);
    }

    if (this.creatingPromise) {
      return this.creatingPromise;
    }
    this.creatingPromise = new Promise((resolve) => {
      UiViewFactory.createUi(R.Prefabs.MailPane, LayerIdx.layer4, (ui: MailPaneController) => {
        this.instance = ui;
        this.creatingPromise = null;
        resolve(ui);
      });
    });
    return this.creatingPromise;
  }
}
