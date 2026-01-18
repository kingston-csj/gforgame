import GameContext from "../../GameContext";
import RewardVo from "../../net/protocol/items/RewardVo";
import { ReqHeroRecruit } from "../../net/protocol/ReqHeroRecruit";
import { ResHeroRecruit } from "../../net/protocol/ResHeroRecruit";

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
      GameContext.instance.WebSocketClient.sendMessage(
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
