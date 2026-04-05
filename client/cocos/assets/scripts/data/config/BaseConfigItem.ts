export default class BaseConfigItem {
  private _id: number;

  private _name: string;

  private _desc: string;

  public constructor(data: any) {
    this._id = data.id;
    this._name = data.name;
    this._desc = data.desc;
  }

  public get id(): number {
    return this._id;
  }

  public get name(): string {
    return this._name;
  }

  public get desc(): string {
    return this._desc;
  }
}
