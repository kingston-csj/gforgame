import { ItemInfo } from "./items/ItemInfo";

export class PushItemChanged {
  public static cmd: number = 4003;
  // item, rune,card 等道具类型
  public type: string;
  public changed: ItemInfo[];
}
