import { _decorator, EditBox } from 'cc';
import GameContext from '../../GameContext';
import ReqGmAction from '../../net/protocol/ReqGmAction';
import ResGmAction from '../../net/protocol/ResGmAction';
import { BaseUiView } from '../../ui/BaseUiView';

const { ccclass, property } = _decorator;

@ccclass('GmPaneView')
export class GmPaneView extends BaseUiView {
  @property(EditBox)
  public itemIdBox: EditBox = null!;

  @property(EditBox)
  public itemNumBox: EditBox = null!;

  @property(EditBox)
  public goldBox: EditBox = null!;

  @property(EditBox)
  public diamondBox: EditBox = null!;

  public start(): void {}

  public onAddItemBtnClick() {
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
