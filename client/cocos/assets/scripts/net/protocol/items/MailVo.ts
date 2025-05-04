import RewardInfo from './RewardInfo';

export class MailVo {
  public id: number;
  public title: string;
  public content: string;
  public time: number;
  public status: number;
  public rewards: RewardInfo[];
}
