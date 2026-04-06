import { ConfigContext } from "../../data/config/container/ConfigContext";
import { BaseModel } from "../../frame/mvc/BaseModel";
import QuestVo from "../../net/protocol/items/QuestVo";
import GameConstants from "../constants/GameConstants";
import GameEvent from "../constants/GameEvent";
export class QuestBoxModel extends BaseModel {
  private static _instance: QuestBoxModel = new QuestBoxModel();

  private quests: Map<number, QuestVo> = new Map();

  // 每日任务奖励积分
  public dailyScore: number = 0;

  // 每日任务奖励索引
  public dailyRewardIndex: number = 0;

  public static get instance(): QuestBoxModel {
    return this._instance;
  }

  public refreshQuest(questVo: QuestVo) {
    this.quests.set(questVo.id, questVo);
  }

  public getQuest(questId: number): QuestVo {
    return this.quests.get(questId);
  }

  public getMainQuest(): QuestVo {
    for (let quest of this.quests.values()) {
        let questData = ConfigContext.configQuestContainer.getRecord(quest.id);
      if (questData.category == GameConstants.Quest.Category.MAIN) {
        return quest;
      }
    }
    return null;
  }

  public getQuestsByCategory(category: number): QuestVo[] {
    let quests: QuestVo[] = [];
    for (let quest of this.quests.values()) {
      let questData = ConfigContext.configQuestContainer.getRecord(quest.id);
      if (questData.category == category) {
        quests.push(quest);
      }
    }
    return quests;
  }
}