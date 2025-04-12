import ConfigContainer from '../ConfigContainer';

import HerostageData from '../model/HerostageData';

export default class ConfigHeroStageContainer extends ConfigContainer<HerostageData> {
  private stageMapper: Map<number, HerostageData> = new Map();

  public constructor() {
    super(HerostageData, HerostageData.fileName);
  }

  public afterLoad() {
    super.afterLoad();
    this.stageMapper.clear();
    for (const item of this._datas.values()) {
      this.stageMapper.set(item.stage, item);
    }
  }

  public getRecordByStage(stage: number): HerostageData | undefined {
    return this.stageMapper.get(stage);
  }
}
