export class MessageDispatch {
  // 绑定cmd与对应的handler
  private static handlers: Map<number, Function> = new Map();

  /**
   * 注册消息处理器
   */
  public static register(cmd: number, handler: Function): void {
    this.handlers.set(cmd, handler);
  }

  public static dispatch(cmd: number, msg: any): void {
    const handler = this.handlers.get(cmd);
    if (handler) {
      handler(msg);
    }
  }
}
