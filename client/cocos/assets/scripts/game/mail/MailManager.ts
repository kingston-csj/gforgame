import { RedDotManager } from '../../frame/reddot/RedDotManager';
import { MailBoxModel } from './MailBoxModel';

export class MailManager {
  private static instance: MailManager;

  public static getInstance(): MailManager {
    if (!MailManager.instance) {
      MailManager.instance = new MailManager();
    }
    return MailManager.instance;
  }

  public refreshRedDots() {
    let rewardAllBtnRedDot = false;
    // 邮件单项红点
    MailBoxModel.getInstance()
      .getMails()
      .forEach((mail) => {
        let hasReward = mail.hasNotReceivedRewards();
        RedDotManager.instance.updateScore(`mail/${mail.id}`, hasReward ? 1 : 0);
        if (hasReward) {
          rewardAllBtnRedDot = true;
        }
      });
    // 一键领奖红点
    RedDotManager.instance.updateScore(`mail/all`, rewardAllBtnRedDot ? 1 : 0);
  }
}
