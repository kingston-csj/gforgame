import Bagpack from './data/user/Bagpack';
import GameContext from './GameContext';
import ResBackpackInfo from './net/ResBackpackInfo';

export class MessageDispatch {
  // 绑定cmd与对应的handler
  private static handlers: Map<number, Function> = new Map();

  public static init(): void {
    MessageDispatch.register(ResBackpackInfo.cmd, (msg: ResBackpackInfo) => {
      GameContext.instance.playerData.Bagpack = new Bagpack(
        new Map(msg.items.map((item) => [item.id, item]))
      );
    });
  }

  /**
   * 注册消息处理器
   */
  private static register(cmd: number, handler: Function): void {
    this.handlers.set(cmd, handler);
  }

  public static dispatch(cmd: number, msg: any): void {
    const handler = this.handlers.get(cmd);
    if (handler) {
      handler(msg);
    }
  }
}
