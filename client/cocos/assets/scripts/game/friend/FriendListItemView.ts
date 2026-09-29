import { _decorator, instantiate, Label, Node, Prefab } from 'cc';
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import { RedDotComponent } from '../../frame/reddot/RedDotCompoent';
import { RedDotManager } from '../../frame/reddot/RedDotManager';
import { TimeUtils } from '../../utils/TimeUtils';
import { RewardItem } from '../reward/RewardItem';
import { FriendVo } from '../../net/protocol/friend/FriendVo';
const { ccclass, property } = _decorator;

@ccclass('FriendListItemView')
export class FriendListItemView extends BaseUiView {
  @property(Label)
  nameTxt: Label;

  @property(Label)
  levelTxt: Label;


  protected start(): void {
  }


  public fillData(friend: FriendVo): void {
    this.nameTxt.string = friend.name;
  }
}
