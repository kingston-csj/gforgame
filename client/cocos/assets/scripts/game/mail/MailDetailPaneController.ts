import { _decorator } from 'cc';
import { BaseController } from '../../frame/mvc/BaseController';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
import { MailDetailView } from './MailDetailView';

const { ccclass, property } = _decorator;

@ccclass('MailDetailPaneController')
export class MailDetailPaneController extends BaseController {
  private static instance: MailDetailPaneController;
  private static creatingPromise: Promise<MailDetailPaneController> | null = null;

  @property(MailDetailView)
  public view: MailDetailView;

  public static openUi(mailId: number) {
    this.getInstance().then((controller) => {
      controller.view.selectedMailId = mailId;
      controller.view.display();
    });
  }

  protected start(): void {
    this.initView(this.view);
  }

  private static getInstance(): Promise<MailDetailPaneController> {
    if (this.instance) {
      return Promise.resolve(this.instance);
    }

    if (this.creatingPromise) {
      return this.creatingPromise;
    }
    this.creatingPromise = new Promise((resolve) => {
      UiViewFactory.createUi(
        R.Prefabs.MailDetailPane,
        LayerIdx.layer4,
        (ui: MailDetailPaneController) => {
          this.instance = ui;
          this.creatingPromise = null;
          resolve(ui);
        }
      );
    });
    return this.creatingPromise;
  }
}
