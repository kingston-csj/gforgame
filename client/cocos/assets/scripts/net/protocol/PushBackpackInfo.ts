import { ItemInfo } from "./items/ItemInfo";

export default class PushBackpackInfo {
  public static cmd: number = 4001;

  public items: ItemInfo[];
}
