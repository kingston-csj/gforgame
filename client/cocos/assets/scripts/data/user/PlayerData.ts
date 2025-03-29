import BagpackModel from '../../game/item/BagpackModel';

export default class PlayerData {
  private _id: string;
  private _name: string;

  private _bagpack: BagpackModel;

  public get Id(): string {
    return this._id;
  }

  public get Name(): string {
    return this._name;
  }

  public get Bagpack(): BagpackModel {
    return this._bagpack;
  }

  public set Bagpack(bagpack: BagpackModel) {
    this._bagpack = bagpack;
  }
}
