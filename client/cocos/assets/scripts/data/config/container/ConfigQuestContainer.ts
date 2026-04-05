import ConfigContainer from '../ConfigContainer';
import QuestData from '../model/QuestData';

export default class ConfigQuestContainer extends ConfigContainer<QuestData> {
  public constructor() {
    super(QuestData, QuestData.fileName);
  }   
}
