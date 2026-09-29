import GameContext from "../../GameContext";
import { RewardVo } from "../../net/protocol/common/RewardVo";
import { ReqHeroRecruit } from "../../net/protocol/hero/ReqHeroRecruit";
import { ResHeroRecruit } from "../../net/protocol/hero/ResHeroRecruit";

export class RecruitSettleModel {
  public static instance: RecruitSettleModel;

  private rewardItems: RewardVo[] = [];

  public static getInstance(): RecruitSettleModel {
    if (!RecruitSettleModel.instance) {
      RecruitSettleModel.instance = new RecruitSettleModel();
    }
    return RecruitSettleModel.instance;
  }

  public setRewardItems(rewardItems: RewardVo[]) {
    this.rewardItems = rewardItems;
  }

  public getRewardItems(): RewardVo[] {
    return this.rewardItems;
  }

  public doRecruit(times: number): Promise<ResHeroRecruit> {
    const promise = new Promise<ResHeroRecruit>((resolve, reject) => {
      GameContext.wsClient.sendMessage(
        ReqHeroRecruit.cmd,
        { counter: times },
        (msg: ResHeroRecruit) => {
          resolve(msg);
        },
      );
    });
    return promise;
  }
}
