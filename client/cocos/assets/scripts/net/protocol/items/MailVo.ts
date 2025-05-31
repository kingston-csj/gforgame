import { MailBoxModel } from '../../../game/mail/MailBoxModel';
import RewardInfo from './RewardInfo';

export class MailVo {
  public id: number;
  public title: string;
  public content: string;
  public time: number;
  public status: number;
  public rewards: RewardInfo[];

  public hasNotReceivedRewards(): boolean {
    return this.status != MailBoxModel.STATUS_RECEIVED && this.rewards.length > 0;
  }
}
