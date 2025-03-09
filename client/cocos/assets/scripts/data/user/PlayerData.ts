export default class PlayerData {
  private _id: string;
  private _name: string;
  private _level: number;
  private _exp: number;

  public get Id(): string {
    return this._id;
  }

  public get Name(): string {
    return this._name;
  }
}
