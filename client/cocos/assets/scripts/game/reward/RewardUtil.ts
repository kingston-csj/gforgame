import RewardVo from "../../net/protocol/items/RewardVo";

export class RewardUtil {
    
  public static parseReward(rewards: string): RewardVo[] {
    let rewardItems: RewardVo[] = [];
    let rewardItemsStr = rewards.split(',');
    for (let rewardItemStr of rewardItemsStr) {
      let rewardItem = rewardItemStr.split('_');
      rewardItems.push({
        type: rewardItem[0],
        value: rewardItem[1],
      });
    }
    return rewardItems;
  }
}