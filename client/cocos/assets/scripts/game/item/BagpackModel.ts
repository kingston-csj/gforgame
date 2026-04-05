import { BaseModel } from "../../frame/mvc/BaseModel";
import { ItemInfo } from "../../net/protocol/items/ItemInfo";

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

  public changeItemByModelId(items: ItemInfo[]): void {
    let exist = false;
    for (let item of items) {
      if (this._items.has(item.uid)) {
        let prev = this._items.get(item.uid);
        prev.count += item.count;
        if (item.count <= 0) {
          this._items.delete(item.uid);
        }
      }else {
        let item = new Item();
        item.cf_id = item.cf_id;
        item.uid = item.uid;
        item.count = item.count;
        item.level = item.level;
        item.extra = item.extra;
        this._items.set(item.uid, item);
      }
    }
    if (!exist) {
      let item = new Item();

      this._items.set(item.uid, item);
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
