import { BaseModel } from '../../ui/BaseModel';

// 玩家基本信息
export default class PlayerData extends BaseModel {
  private static _instance: PlayerData = new PlayerData();

  public static get instance(): PlayerData {
    return this._instance;
  }

  private _id: string;
  private _name: string;
  private _fighting: number;

  public get Id(): string {
    return this._id;
  }

  public get Name(): string {
    return this._name;
  }

  public get Fighting(): number {
    return this._fighting;
  }

  public set Fighting(value: number) {
    this._fighting = value;
  }
}
