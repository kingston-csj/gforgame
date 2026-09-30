import { _decorator, instantiate, Label, Node, Prefab } from 'cc';
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import { FriendItem } from './FriendModel';
const { ccclass, property } = _decorator;

@ccclass('FriendSearchItemView')
export class FriendSearchItemView extends BaseUiView {
  @property(Label)
  nameTxt: Label;

  @property(Label)
  levelTxt: Label;

  @property(Node)
  addBtn: Node;

  protected start(): void {
  }


  public fillData(friend: FriendItem): void {
    this.nameTxt.string = friend.name;
  }
}
