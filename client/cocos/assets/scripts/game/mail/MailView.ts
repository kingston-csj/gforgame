import { _decorator, instantiate, Node, Prefab } from 'cc';
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import { RedDotComponent } from '../../frame/reddot/RedDotCompoent';
import { RedDotManager } from '../../frame/reddot/RedDotManager';
import GameContext from '../../GameContext';
import { ReqMailDeleteAll } from '../../net/protocol/ReqMailDeleteAll';
import { ReqMailGetAllReward } from '../../net/protocol/ReqMailGetAllReward';
import { ResMailDeleteAll } from '../../net/protocol/ResMailDeleteAll';
import { ResMailGetAllReward } from '../../net/protocol/ResMailGetAllReward';
import { MailBoxModel } from './MailBoxModel';
import { MailItemView } from './MailItemView';
import { MailManager } from './MailManager';
const { ccclass, property } = _decorator;

@ccclass('MailView')
export class MailView extends BaseUiView {
  @property(Node)
  mailContainer: Node;

  @property(Prefab)
  mailItemPrefab: Prefab;

  // 一键领奖
  @property(Node)
  rewardBtn: Node;

  @property(Node)
  rewardRedDot: Node;

  // 一键删除
  @property(Node)
  deleteBtn: Node;

  @property(Node)
  closeBtn: Node;

  protected start(): void {
    this.registerClickEvent(this.rewardBtn, this.onRewardBtnClick, this);
    this.registerClickEvent(this.deleteBtn, this.onDeleteBtnClick, this);
    this.registerClickEvent(this.closeBtn, this.hide, this);
  }

  private onRewardBtnClick(): void {
    GameContext.instance.WebSocketClient.sendMessage(
      ReqMailGetAllReward.cmd,
      new ReqMailGetAllReward(),
      (res: ResMailGetAllReward) => {
        if (res.code === 0) {
          // 所有邮件，设置为已领奖
          MailBoxModel.getInstance()
            .getMails()
            .forEach((mail) => {
              mail.status = MailBoxModel.STATUS_RECEIVED;
            });

          MailManager.getInstance().refreshRedDots();
        }
      }
    );
  }

  private onDeleteBtnClick(): void {
    GameContext.instance.WebSocketClient.sendMessage(
      ReqMailDeleteAll.cmd,
      new ReqMailDeleteAll(),
      (res: ResMailDeleteAll) => {
        if (res.removed.length > 0) {
          // 删除邮件
          MailBoxModel.getInstance().deleteMails(res.removed);
          // 刷新邮件列表
          this.onDisplay();
        }
      }
    );
  }

  protected onDisplay(): void {
    this.mailContainer.children.forEach((child) => {
      child.destroy();
    });
    const mails = MailBoxModel.getInstance().getMails();
    for (const mail of mails) {
      const mailItem = instantiate(this.mailItemPrefab);
      mailItem.setParent(this.mailContainer);
      mailItem.getComponent(MailItemView).fillData(mail);
    }

    // 绑定红点
    RedDotManager.instance.binding(`mail/all`, this.rewardRedDot.getComponent(RedDotComponent));
  }
}
