import { _decorator, Toggle, Node, EditBox, Button, director } from 'cc';

import UiView from '../../ui/UiView';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import { UIViewController } from '../../ui/UiViewController';
import { BagPanelController } from '../item/BagPanelController';

const { ccclass, property } = _decorator;

@ccclass('MainPaneController')
export class MainPaneController extends UIViewController {
  @property(Node)
  bagPane: Node;

  private static instance: MainPaneController;

  public static display() {
    if (MainPaneController.instance) {
      MainPaneController.instance.display();
    } else {
      MainPaneController.instance = new MainPaneController();

      UiView.createUi(R.mainPane, LayerIdx.layer1, (ui: MainPaneController) => {
        MainPaneController.instance = ui;
        ui.display();
      });
    }
  }

  protected start(): void {
    this.registerClickEvent(this.bagPane, this.onBagClick, this);
  }

  onBagClick() {
    BagPanelController.display();
  }
}
