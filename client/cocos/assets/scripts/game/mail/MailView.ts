import { _decorator, instantiate, Node, Prefab } from 'cc';
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import { RedDotComponent } from '../../frame/reddot/RedDotCompoent';
import { RedDotManager } from '../../frame/reddot/RedDotManager';
import GameContext from '../../GameContext';
import { ReqMailDeleteAll } from '../../net/protocol/mail/ReqMailDeleteAll';
import { ReqMailGetAllRewards } from '../../net/protocol/mail/ReqMailGetAllRewards';
import { ResMailDeleteAll } from '../../net/protocol/mail/ResMailDeleteAll';
import { ResMailGetAllRewards } from '../../net/protocol/mail/ResMailGetAllRewards';
import { MailBoxModel } from './MailBoxModel';
import { MailItemView } from './MailItemView';
import { MailManager } from './MailManager';
const { ccclass, property } = _decorator;

@ccclass('MailView')
export class MailView extends BaseUiView {
  @property(Node)
  mailContainer: Node;
  @property(Node)
  closeBtn: Node;
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


  protected start(): void {
    this.registerClickEvent(this.rewardBtn, this.onRewardBtnClick, this);
    this.registerClickEvent(this.deleteBtn, this.onDeleteBtnClick, this);
    this.registerClickEvent(this.closeBtn, this.hide, this);
  }

  private onRewardBtnClick(): void {
    GameContext.wsClient.sendMessage(
      ReqMailGetAllRewards.cmd,
      new ReqMailGetAllRewards(),
      (res: ResMailGetAllRewards) => {
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
    GameContext.wsClient.sendMessage(
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
