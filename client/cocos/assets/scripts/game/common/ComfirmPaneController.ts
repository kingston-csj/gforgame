import { _decorator, Button } from 'cc';
import { BaseController } from '../../ui/BaseController';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
import { ComfirmView } from './ComfirmView';
const { ccclass, property } = _decorator;

// 二次弹窗确认框
@ccclass('ComfirmPaneController')
export class ComfirmPaneController extends BaseController {
  @property(ComfirmView)
  comfirmView: ComfirmView | null = null;

  private static instance: ComfirmPaneController;

  private static creatingPromise: Promise<ComfirmPaneController> | null = null;

  start() {
    this.initView(this.comfirmView);
  }

  public static show(title: string, content: string, confirmCallback: () => void) {
    this.getInstance().then((controller) => {
      if (controller.comfirmView) {
        controller.comfirmView.display();
        controller.comfirmView.setTitle(title);
        controller.comfirmView.setContent(content);
        controller.comfirmView.confirmButton.node.on(
          Button.EventType.CLICK,
          confirmCallback,
          controller
        );
        controller.comfirmView.cancelButton.node.on(
          Button.EventType.CLICK,
          controller.comfirmView.hide,
          controller
        );
      }
    });
  }

  public static hide() {
    this.getInstance().then((controller) => {
      if (controller.comfirmView) {
        controller.comfirmView.hide();
      }
    });
  }

  private static getInstance(): Promise<ComfirmPaneController> {
    if (this.instance) {
      return Promise.resolve(this.instance);
    }
    if (this.creatingPromise) {
      return this.creatingPromise;
    }
    this.creatingPromise = new Promise((resolve) => {
      UiViewFactory.createUi(
        R.Prefabs.ComfirmPane,
        LayerIdx.layer5,
        (ui: ComfirmPaneController) => {
          this.instance = ui;
          this.creatingPromise = null;
          resolve(ui);
        }
      );
    });
    return this.creatingPromise;
  }
}
