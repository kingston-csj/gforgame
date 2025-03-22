export default class Bagpack {
  private _items: Map<number, Item> = new Map();

  constructor(items: Map<number, Item>) {
    this._items = items;
  }

  public getItems(): Array<Item> {
    return Array.from(this._items.values());
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
