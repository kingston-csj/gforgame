import BaseConfigItem from '../BaseConfigItem';

export default class QuestData extends BaseConfigItem {
  public static fileName = 'quest';

  private _category: number;

  private _target: string;

  private _rewards: string;

  public get rewards(): string {
    return this._rewards;
  }

  public get category(): number {
    return this._category;
  }

  public get target(): string {
    return this._target;
  }

  public get type(): string {
    return this._type;
  }

  private _type: string;

  public constructor(data:any) {
            super(data);
            this._category = data.category;
            this._type = data.type;
            this._target = data.target;
            this._rewards = data.rewards;
  }

}
