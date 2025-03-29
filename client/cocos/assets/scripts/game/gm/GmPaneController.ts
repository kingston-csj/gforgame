import { _decorator, Component, Node, EditBox, Button, director } from 'cc';
import { BaseUiView } from '../../ui/BaseUiView';
import UiViewFactory from '../../ui/UiViewFactory';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import GameContext from '../../GameContext';
import ReqGmAction from '../../net/ReqGmAction';
import ResGmAction from '../../net/ResGmAction';

const { ccclass, property } = _decorator;

@ccclass('GmPaneController')
export class GmPaneController extends BaseUiView {
  private static instance: GmPaneController;

  @property(Node)
  container: Node = null!;

  @property(EditBox)
  public itemIdBox: EditBox = null!;

  @property(EditBox)
  public itemNumBox: EditBox = null!;

  public static display() {
    if (GmPaneController.instance) {
      GmPaneController.instance.display();
    } else {
      GmPaneController.instance = new GmPaneController();

      UiViewFactory.createUi(R.gmPane, LayerIdx.layer2, () => {});
    }
  }

  public onShowBtnToggle(event: Event, customEventData: string) {
    this.container.active = !this.container.active;
  }

  public onAddBtnClick() {
    const itemId = this.itemIdBox.string;
    const itemNum = this.itemNumBox.string;

    GameContext.instance.WebSocketClient.sendMessage(
      ReqGmAction.cmd,
      {
        topic: 'add_item',
        params: itemId + '=' + itemNum,
      },
      (msg: ResGmAction) => {}
    );
  }
}
