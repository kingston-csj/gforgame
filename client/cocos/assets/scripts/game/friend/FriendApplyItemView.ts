import { _decorator, instantiate, Label, Node, Prefab } from 'cc';
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import { FriendApplyVo } from '../../net/protocol/friend/FriendApplyVo';
const { ccclass, property } = _decorator;

@ccclass('FriendApplyItemView')
export class FriendApplyItemView extends BaseUiView {
  @property(Label)
  nameTxt: Label;

  @property(Label)
  levelTxt: Label;

  @property(Node)
  addBtn: Node;

  @property(Node)
  rejectBtn: Node;

  protected start(): void {
  }


  public fillData(friend: FriendApplyVo): void {
    this.nameTxt.string = friend.fromName;
  }
}
