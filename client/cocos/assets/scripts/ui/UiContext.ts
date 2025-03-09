import { Node } from 'cc';

export default class UiContext {
  private static _layers: Array<Node>;

  public static init(layer1: Node, layer2: Node, layer3: Node, layer4: Node, layer5: Node) {
    this._layers = [layer1, layer2, layer3, layer4, layer5];
  }

  public static getLayer(layerIdx: number) {
    return this._layers[layerIdx];
  }
}
