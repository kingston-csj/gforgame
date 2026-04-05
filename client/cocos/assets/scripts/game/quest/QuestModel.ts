import { ConfigContext } from "../../data/config/container/ConfigContext";
import { BaseModel } from "../../frame/mvc/BaseModel";
import QuestVo from "../../net/protocol/items/QuestVo";
import GameConstants from "../constants/GameConstants";
import GameEvent from "../constants/GameEvent";
export class QuestModel extends BaseModel {
  private static _instance: QuestModel = new QuestModel();

  private quests: Map<number, QuestVo> = new Map();

  public static get instance(): QuestModel {
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
}