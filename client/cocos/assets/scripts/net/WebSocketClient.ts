import { MessageDispatch } from '../MessageDispatch';

export class WebSocketClient {
  private ws: WebSocket | null = null;

  private _index: number = 0;
  private handles: Map<number, Function> = new Map();

  /**
   * 向服务器发送消息
   * @param msg 消息
   */
  public sendMessage(msgId: number, msg: any, callback: Function): void {
    this._index++;

    if (callback) {
      this.handles.set(this._index, callback);
    }

    this.sendBytes(msgId, this._index, msg);
  }

  /**
   * 消息发送（二进制）
   */
  private sendBytes(msgId: number, index: number, msg: any): boolean {
    let headerSize = 12;
    let json = JSON.stringify(msg);
    let msgSize = json.length;
    let buffer = new Uint8Array(headerSize + msgSize);
    // 将命令和长度转换为字节并复制到缓冲区
    const cmdBytes = this.intToBytes(msgId);
    // 流水号统一为0
    const indexBytes = this.intToBytes(index);
    const lenBytes = this.intToBytes(msgSize);

    buffer.set(cmdBytes, 0);
    buffer.set(indexBytes, 4);
    buffer.set(lenBytes, 8);
    buffer.set(new TextEncoder().encode(json), headerSize);
    this.ws.send(buffer);
    return true;
  }

  private intToBytes(n) {
    const buf = new Uint8Array(4);
    buf[0] = (n >> 24) & 0xff;
    buf[1] = (n >> 16) & 0xff;
    buf[2] = (n >> 8) & 0xff;
    buf[3] = n & 0xff;
    return buf;
  }

  public fillCallback(index: number, msg: any): boolean {
    const callback = this.handles.get(index);
    if (callback) {
      callback(msg);
      return true;
    }
    return false;
  }

  /**
   * 连接服务器
   * @param url 服务器地址(ws://localhost:9527/ws)
   */
  public connect(url: string): void {
    this.ws = new WebSocket(url);
    this.ws.onopen = (evt: Event) => {
      console.info('建立连接');
      this.onConnect(evt);
    };
    this.ws.onmessage = (evt: MessageEvent) => {
      this.onMessage(evt);
    };
    this.ws.onclose = (evt: CloseEvent) => {
      this.onClose(evt);
    };
  }

  protected onConnect(evt: Event): void {}

  protected onMessage(evt: MessageEvent): void {
    // 检查是否为二进制数据
    if (evt.data instanceof Blob) {
      evt.data.arrayBuffer().then((buffer) => {
        const uint8Array = new Uint8Array(buffer);
        // 解析头部信息 (12字节)
        const headerSize = 12;
        const cmd = this.bytesToInt(uint8Array.slice(0, 4));
        const index = this.bytesToInt(uint8Array.slice(4, 8));
        const msgSize = this.bytesToInt(uint8Array.slice(8, 12));

        // 解析消息体
        const msgBytes = uint8Array.slice(headerSize, headerSize + msgSize);
        const decoder = new TextDecoder();
        const msgStr = decoder.decode(msgBytes);
        const body = JSON.parse(msgStr);

        const callback = this.handles.get(index);
        if (callback) {
          // 属于客户端回调
          callback(body);
          return;
        } else {
          // 属于服务器主动推送
          MessageDispatch.dispatch(cmd, body);
        }
      });
    } else {
      // 兼容原有的JSON格式
      let frame = JSON.parse(evt.data);
      let cmd = frame.cmd;
      let index = frame.index;
      let body = JSON.parse(frame.msg);

      const callback = this.handles.get(index);
      if (callback) {
        callback(body);
        return;
      } else {
        MessageDispatch.dispatch(cmd, body);
      }
    }
  }

  private bytesToInt(bytes: Uint8Array): number {
    return (bytes[0] << 24) | (bytes[1] << 16) | (bytes[2] << 8) | bytes[3];
  }

  protected onClose(evt: CloseEvent): void {}

  protected onError(evt: ErrorEvent): void {}
}
