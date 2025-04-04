import { _decorator, EditBox, Node } from 'cc';
import GameContext from '../../GameContext';
import ReqGmAction from '../../net/ReqGmAction';
import ResGmAction from '../../net/ResGmAction';
import { BaseUiView } from '../../ui/BaseUiView';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';

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

  @property(EditBox)
  public goldBox: EditBox = null!;

  @property(EditBox)
  public diamondBox: EditBox = null!;

  public static display() {
    if (GmPaneController.instance) {
      GmPaneController.instance.display();
    } else {
      GmPaneController.instance = new GmPaneController();

      UiViewFactory.createUi(R.Prefabs.GmPane, LayerIdx.layer2, () => {});
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

  public onAddGoldBtnClick() {
    const gold = this.goldBox.string;

    GameContext.instance.WebSocketClient.sendMessage(
      ReqGmAction.cmd,
      {
        topic: 'add_gold',
        params: gold,
      },
      (msg: ResGmAction) => {}
    );
  }

  public onAddDiamondBtnClick() {
    const diamond = this.diamondBox.string;

    GameContext.instance.WebSocketClient.sendMessage(
      ReqGmAction.cmd,
      {
        topic: 'add_diamond',
        params: diamond,
      },
      (msg: ResGmAction) => {}
    );
  }
}
