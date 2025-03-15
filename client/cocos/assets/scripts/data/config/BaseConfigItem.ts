export default class BaseConfigItem {
  private _id: number;

  public constructor(data: any) {
    this._id = data.id;
  }

  public get id(): number {
    return this._id;
  }
}
