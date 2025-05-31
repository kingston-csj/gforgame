import { _decorator, Button, instantiate, Label, Node, Prefab } from 'cc';

import GameContext from '../../GameContext';

import { BaseUiView } from '../../frame/mvc/BaseUiView';
import { RedDotManager } from '../../frame/reddot/RedDotManager';
import { ReqMailGetReward } from '../../net/protocol/ReqMailGetReward';
import { ResMailGetReward } from '../../net/protocol/ResMailGetReward';
import { TimeUtils } from '../../utils/TimeUtils';
import { RewardItem } from '../reward/RewardItem';
import { MailBoxModel } from './MailBoxModel';
import { MailManager } from './MailManager';
const { ccclass, property } = _decorator;

@ccclass('MailDetailView')
export class MailDetailView extends BaseUiView {
  @property(Node)
  rewardContainer: Node;

  @property(Prefab)
  rewardItemPrefab: Prefab;

  @property(Label)
  titleLabel: Label;

  @property(Label)
  contentLabel: Label;

  @property(Label)
  timeLabel: Label;

  // 一键领奖
  @property(Node)
  rewardBtn: Node;

  @property(Node)
  closeBtn: Node;

  selectedMailId: number = 0;

  protected start(): void {
    this.registerClickEvent(this.rewardBtn, this.onRewardBtnClick, this);
    this.registerClickEvent(this.closeBtn, this.hide, this);
  }

  private onRewardBtnClick(): void {
    let req = new ReqMailGetReward();
    req.id = this.selectedMailId;
    GameContext.instance.WebSocketClient.sendMessage(
      ReqMailGetReward.cmd,
      req,
      (res: ResMailGetReward) => {
        if (res.code === 0) {
          // 领奖成功
          MailBoxModel.getInstance().getMail(this.selectedMailId).status =
            MailBoxModel.STATUS_RECEIVED;
          this.rewardBtn.getComponent(Button).interactable = false;
          RedDotManager.instance.updateScore(`mail/${this.selectedMailId}`, 0);
          // 刷新红点
          MailManager.getInstance().refreshRedDots();
        }
      }
    );
  }

  protected onDisplay(): void {
    this.rewardContainer.children.forEach((child) => {
      child.destroy();
    });
    if (this.selectedMailId > 0) {
      const mail = MailBoxModel.getInstance().getMail(this.selectedMailId);
      this.titleLabel.string = mail.title;
      this.contentLabel.string = mail.content;
      let expiredTime = mail.time + 30 * TimeUtils.ONE_DAY;
      this.timeLabel.string = TimeUtils.getLeftTimeTips(expiredTime);
      for (const reward of mail.rewards) {
        const rewardItem = instantiate(this.rewardItemPrefab);
        rewardItem.setParent(this.rewardContainer);
        rewardItem.getComponent(RewardItem).fillData(reward);
      }

      this.rewardBtn.getComponent(Button).interactable =
        mail.status != MailBoxModel.STATUS_RECEIVED;
    }
  }
}
