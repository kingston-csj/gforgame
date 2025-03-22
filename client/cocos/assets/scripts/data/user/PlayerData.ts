import Bagpack from './Bagpack';

export default class PlayerData {
  private _id: string;
  private _name: string;

  private _bagpack: Bagpack;

  public get Id(): string {
    return this._id;
  }

  public get Name(): string {
    return this._name;
  }

  public get Bagpack(): Bagpack {
    return this._bagpack;
  }

  public set Bagpack(bagpack: Bagpack) {
    this._bagpack = bagpack;
  }
}
