import { Attribute } from './items/Attribute';

export default class PushHeroAttrChanged {
  public static cmd: number = 5007;
  public heroId: number;
  public attrs: Attribute[];
  public fight: number;
}
