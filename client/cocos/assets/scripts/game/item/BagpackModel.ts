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

  public addItem(item: Item): void {
    this._items.set(item.id, item);
  }

  public removeItem(item: Item): void {
    this._items.delete(item.id);
  }

  public getItemCount(itemId: number): number {
    return this._items.get(itemId)?.count || 0;
  }
}

export class Item {
  public id: number;
  public count: number;
}
