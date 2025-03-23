import { Item } from '../data/user/Bagpack';
import { RewardInfo } from '../net/ResHeroRecruit';

export class RecruitSettleModel {
  public static instance: RecruitSettleModel;

  private rewardItems: RewardInfo[] = [];

  public static getInstance(): RecruitSettleModel {
    if (!RecruitSettleModel.instance) {
      RecruitSettleModel.instance = new RecruitSettleModel();
    }
    return RecruitSettleModel.instance;
  }

  public setRewardItems(rewardItems: RewardInfo[]) {
    this.rewardItems = rewardItems;
  }

  public getRewardItems(): RewardInfo[] {
    return this.rewardItems;
  }
}
