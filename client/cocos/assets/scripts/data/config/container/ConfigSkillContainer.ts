import ConfigContainer from '../ConfigContainer';
import SkillData from '../model/SkillData';

export default class ConfigSkillContainer extends ConfigContainer<SkillData> {
  private skillMapper: Map<number, SkillData>;

  public constructor() {
    super(SkillData, SkillData.fileName);
  }

  public afterLoad() {
    this.skillMapper = new Map();
    for (const skill of this._datas.values()) {
      this.skillMapper.set(skill.skillId, skill);
    }
  }

  public getSkill(skillId: number): SkillData | undefined {
    return this.skillMapper?.get(skillId);
  }
}
