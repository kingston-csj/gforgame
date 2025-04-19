import { Attribute } from '../../net/protocol/items/Attribute';
import { AttributeTypes } from './attributetypes';

export class AttributeBox {
  private attrs: Map<string, number> = new Map();

  public constructor(attrs: Attribute[]) {
    for (const attr of attrs) {
      this.attrs.set(attr.attrType, attr.value);
    }
  }

  public getAttr(attrType: string): number {
    return this.attrs.get(attrType);
  }

  public getAttack(): number {
    return this.attrs.get(AttributeTypes.Attack);
  }

  public getDefense(): number {
    return this.attrs.get(AttributeTypes.Defense);
  }

  public getSpeed(): number {
    return this.attrs.get(AttributeTypes.Speed);
  }

  public getHp(): number {
    return this.attrs.get(AttributeTypes.Hp);
  }
}
