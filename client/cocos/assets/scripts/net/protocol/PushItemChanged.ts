export class PushItemChanged {
  public static cmd: number = 4003;

  /// <summary>
  /// 物品ID
  /// </summary>
  public itemId: number;
  /// <summary>
  /// 物品数量
  /// </summary>
  public count: number;
}
