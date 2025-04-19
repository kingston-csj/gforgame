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
  private _camp: number;

  public get Id(): string {
    return this._id;
  }

  public get name(): string {
    return this._name;
  }

  public set name(value: string) {
    this._name = value;
  }

  public get fighting(): number {
    return this._fighting;
  }

  public set fighting(value: number) {
    this._fighting = value;
  }

  public get camp(): number {
    return this._camp;
  }

  public set camp(value: number) {
    this._camp = value;
  }
}
