export default class BagpackModel {
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
    const item = this._items.get(itemId);
    if (item) {
      item.count += count;
    } else {
      this._items.set(itemId, { id: itemId, count });
    }
    if (item.count <= 0) {
      this._items.delete(itemId);
    }
  }
}

export class Item {
  public id: number;
  public count: number;
}
