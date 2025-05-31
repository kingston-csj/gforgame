import { BaseModel } from '../../frame/mvc/BaseModel';

export class PurseModel extends BaseModel {
  private static instance: PurseModel;

  private _diamond: number = 0;
  private _gold: number = 0;

  public static getInstance(): PurseModel {
    if (!PurseModel.instance) {
      PurseModel.instance = new PurseModel();
    }
    return PurseModel.instance;
  }

  // 使用 getter/setter 来触发数据变化通知
  get diamond(): number {
    return this._diamond;
  }

  set diamond(value: number) {
    if (this._diamond !== value) {
      this._diamond = value;
      this.notifyChange('diamond', value);
    }
  }

  get gold(): number {
    return this._gold;
  }

  set gold(value: number) {
    if (this._gold !== value) {
      this._gold = value;
      this.notifyChange('gold', value);
    }
  }
}
