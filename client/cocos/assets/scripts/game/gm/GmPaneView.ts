import { _decorator, EditBox } from "cc";
import { BaseUiView } from "../../frame/mvc/BaseUiView";
import GameContext from "../../GameContext";
import ReqGmAction from "../../net/protocol/ReqGmAction";
import ResGmAction from "../../net/protocol/ResGmAction";

const { ccclass, property } = _decorator;

@ccclass("GmPaneView")
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
        args: "add_items " + itemId + "=" + itemNum,
      },
      (msg: ResGmAction) => {},
    );
  }

  public onAddGoldBtnClick() {
    const gold = this.goldBox.string;

    GameContext.instance.WebSocketClient.sendMessage(
      ReqGmAction.cmd,
      {
        args: "add_gold " + gold,
      },
      (msg: ResGmAction) => {},
    );
  }

  public onAddDiamondBtnClick() {
    const diamond = this.diamondBox.string;

    GameContext.instance.WebSocketClient.sendMessage(
      ReqGmAction.cmd,
      {
        args: "add_diamond " + diamond,
      },
      (msg: ResGmAction) => {},
    );
  }
}
