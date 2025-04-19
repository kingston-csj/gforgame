import { BaseModel } from '../../ui/BaseModel';

export default class BagpackModel extends BaseModel {
  private static instance: BagpackModel;

  public static getInstance(): BagpackModel {
    if (!BagpackModel.instance) {
      BagpackModel.instance = new BagpackModel();
    }
    return BagpackModel.instance;
  }

  private _items: Map<number, Item> = new Map();

  public reset(data: Map<number, Item>) {
    this._items = data;
    this.notifyChange('item', this._items);
  }

  public getItems(): Array<Item> {
    return Array.from(this._items.values());
  }

  public getItemByModelId(itemId: number): Item | undefined {
    return this._items.get(itemId);
  }

  public removeItem(item: Item): void {
    this._items.delete(item.id);
  }

  public getItemCount(itemId: number): number {
    return this._items.get(itemId)?.count || 0;
  }

  public changeItemByModelId(itemId: number, count: number): void {
    let item = this._items.get(itemId);
    if (item) {
      item.count += count;
    } else {
      item = { id: itemId, count };
      this._items.set(itemId, item);
    }
    if (item.count <= 0) {
      this._items.delete(itemId);
    }
    this.notifyChange('item', this._items);
  }
}

export class Item {
  public id: number;
  public count: number;
}
