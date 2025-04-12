import { AttributeBox } from '../../../game/attribute/attributebox';
import { Attribute } from './Attribute';
export class HeroVo {
  public id: number;
  public level: number;
  public position: number;
  public stage: number;
  public attrs: Attribute[];
  public fight: number;

  public attrBox: AttributeBox;
}
