import GameContext from '../../GameContext';
import RewardInfo from '../../net/protocol/items/RewardInfo';
import { ReqHeroRecruit } from '../../net/protocol/ReqHeroRecruit';
import { ResHeroRecruit } from '../../net/protocol/ResHeroRecruit';

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

  public doRecruit(times: number): Promise<ResHeroRecruit> {
    const promise = new Promise<ResHeroRecruit>((resolve, reject) => {
      GameContext.instance.WebSocketClient.sendMessage(
        ReqHeroRecruit.cmd,
        { times },
        (msg: ResHeroRecruit) => {
          resolve(msg);
        }
      );
    });
    return promise;
  }
}
