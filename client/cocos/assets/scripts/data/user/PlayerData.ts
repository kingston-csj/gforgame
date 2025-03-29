// 玩家基本信息
export default class PlayerData {
  private _id: string;
  private _name: string;

  public get Id(): string {
    return this._id;
  }

  public get Name(): string {
    return this._name;
  }
}
