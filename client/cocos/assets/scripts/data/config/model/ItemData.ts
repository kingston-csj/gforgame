import BaseConfigItem from '../BaseConfigItem';

export default class Config_itemData extends BaseConfigItem {
  public static fileName: string = 'itemData';

  private _type: number;
  public get type(): number {
    return this._type;
  }

  private _quality: number;
  public get quality(): number {
    return this._quality;
  }

  private _tips: string;
  public get tips(): string {
    return this._tips;
  }

  private _icon: string;
  public get icon(): string {
    return this._icon;
  }

  public constructor(data: any) {
    super(data);
    this._type = data['type'];
    this._quality = data['quality'];
    this._tips = data['tips'];
    this._icon = data['icon'];
  }
}
