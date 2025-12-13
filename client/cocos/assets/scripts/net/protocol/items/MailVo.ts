import { MailBoxModel } from "../../../game/mail/MailBoxModel";
import RewardVo from "./RewardVo";

export class MailVo {
  public id: number;
  public title: string;
  public content: string;
  public time: number;
  public status: number;
  public rewards: RewardVo[];

  public hasNotReceivedRewards(): boolean {
    return (
      this.status != MailBoxModel.STATUS_RECEIVED && this.rewards.length > 0
    );
  }
}
