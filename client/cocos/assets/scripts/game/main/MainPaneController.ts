import { _decorator, Component, Node, EditBox, Button, director } from 'cc';

import UiView from '../../ui/UiView';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import { UIViewController } from '../../ui/UiViewController';

const { ccclass, property } = _decorator;

@ccclass('MainPaneController')
export class MainPaneController extends UIViewController {
  @property(Button)
  gmButton: Button = null!;

  private static instance: MainPaneController;

  public static display() {
    if (MainPaneController.instance) {
      MainPaneController.instance.display();
    } else {
      MainPaneController.instance = new MainPaneController();

      UiView.createUi(R.mainPane, LayerIdx.layer2, () => {});
    }
  }

  start() {}

  public static close() {
    console.info('------------MainPaneController.close------------');
  }
}
