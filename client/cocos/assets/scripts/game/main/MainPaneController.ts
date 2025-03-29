import { _decorator, Node } from 'cc';

import { BaseUiView } from '../../ui/BaseUiView';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
import { BagPanelController } from '../item/BagPanelController';

import { RecruitPaneController } from '../recruit/RecruitPaneController';

const { ccclass, property } = _decorator;

@ccclass('MainPaneController')
export class MainPaneController extends BaseUiView {
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

      UiViewFactory.createUi(R.Prefabs.MainPane, LayerIdx.layer1, (ui: MainPaneController) => {
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
