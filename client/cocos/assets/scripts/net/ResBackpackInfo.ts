export default class ResBackpackInfo {
  public static cmd: number = 4001;

  public items: ItemInfo[] ; 
}

export class ItemInfo {
  public id: number;
  public count: number;
}
