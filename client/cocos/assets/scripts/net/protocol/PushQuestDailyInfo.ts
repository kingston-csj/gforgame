import QuestVo from "./items/QuestVo";
import { ReqMailGetAllReward } from "./ReqMailGetAllReward";

export default class PushQuestDailyInfo {
    public static cmd = 798;

    public dailyRewardIndex: number = 0;

    public dailyScore: number = 0;

    public quests: QuestVo[] = [];
}