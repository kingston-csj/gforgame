import { _decorator, Toggle, Node, EditBox, Button, director } from 'cc';

import UiView from '../../ui/UiView';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import { UIViewController } from '../../ui/UiViewController';
import { BagPanelController } from '../item/BagPanelController';

import { RecruitPaneController } from '../recruit/RecruitPaneController';

const { ccclass, property } = _decorator;

@ccclass('MainPaneController')
export class MainPaneController extends UIViewController {
  @property(Node)
  bagPane: Node;

  @property(Node)
  recruitePane: Node;

  private static instance: MainPaneController;

  public static openUi() {
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
    this.registerClickEvent(this.recruitePane, this.onRecruiteClick, this);
  }

  onBagClick() {
    BagPanelController.openUi();
  }

  onRecruiteClick() {
    RecruitPaneController.openUi();
  }
}
