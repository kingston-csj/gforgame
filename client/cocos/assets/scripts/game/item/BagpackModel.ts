import { BaseModel } from "../../frame/mvc/BaseModel";

export default class BagpackModel extends BaseModel {
  private static instance: BagpackModel;

  public static getInstance(): BagpackModel {
    if (!BagpackModel.instance) {
      BagpackModel.instance = new BagpackModel();
    }
    return BagpackModel.instance;
  }

  private _items: Map<string, Item> = new Map();

  public reset(data: Map<string, Item>) {
    this._items = data;
    this.notifyChange("item", this._items);
  }

  public getItems(): Array<Item> {
    return Array.from(this._items.values());
  }

  public removeItem(item: Item): void {
    this._items.delete(item.uid);
  }

  public getItemCount(itemId: number): number {
    let sum = 0;
    for (let item of this._items.values()) {
      if (item.cf_id == itemId) {
        sum += item.count;
      }
    }

    return sum;
  }

  public changeItemByModelId(itemId: number, count: number): void {
    for (let item of this._items.values()) {
      if (item.cf_id == itemId) {
        item.count += count;
        if (item.count <= 0) {
          this._items.delete(item.uid);
        }
      }
    }
    this.notifyChange("item", this._items);
  }
}

export class Item {
  public cf_id: number;
  public uid: string;
  public count: number;
  public level: number;
  public extra: string;
}
