/**
 * 对webSocket的封装
 */
(function ($) {
  $.config = {
    url: "", //链接地址
  };

  $.init = function (config) {
    this.config = config;
    return this;
  };

  /**
   * 连接webcocket
   */
  $.connect = function () {
    var protocol = "ws:";
    this.host = protocol + this.config.url;

    window.WebSocket = window.WebSocket || window.MozWebSocket;
    if (!window.WebSocket) {
      // 检测浏览器支持
      this.error("Error: WebSocket is not supported .");
      return;
    }
    this.socket = new WebSocket(this.host); // 创建连接并注册响应函数
    this.socket.onopen = function () {
      $.onopen();
    };
    this.socket.onmessage = function (message) {
      $.onmessage(message);
    };
    this.socket.onclose = function () {
      $.onclose();
      $.socket = null; // 清理
    };
    this.socket.onerror = function (errorMsg) {
      $.onerror(errorMsg);
    };
    return this;
  };

  /**
   * 自定义异常函数
   * @param {Object} errorMsg
   */
  $.error = function (errorMsg) {
    this.onerror(errorMsg);
  };

  /**
   * 消息发送（json）
   */
  $.send = function (msgId, msg) {
    if (this.socket) {
      var req = {
        cmd: msgId,
        msg: JSON.stringify(msg),
      };
      this.socket.send(JSON.stringify(req));
      return true;
    }
    this.error("please connect to the server first !!!");
    return false;
  };

  /**
   * 消息发送（二进制）
   */
  $.sendBytes = function (msgId, msg) {
    if (this.socket) {
      let headerSize = 8
      let json = JSON.stringify(msg);
      let msgSize = json.length;
      let buffer = new Uint8Array(headerSize + msgSize);
      // 将命令和长度转换为字节并复制到缓冲区
      const cmdBytes = intToBytes(msgId);
      const lenBytes = intToBytes(msgSize);

      buffer.set(cmdBytes, 0);
      buffer.set(lenBytes, 4);
      buffer.set(new TextEncoder().encode(json), headerSize);
      this.socket.send(buffer);
      return true;
    }
    this.error("please connect to the server first !!!");
    return false;
  };

  function intToBytes(n) {
    const buf = new Uint8Array(4);
    buf[0] = (n >> 24) & 0xFF;
    buf[1] = (n >> 16) & 0xFF;
    buf[2] = (n >> 8) & 0xFF;
    buf[3] = n & 0xFF;
    return buf;
  }

  $.close = function () {
    if (this.socket !== undefined && this.socket != null) {
      this.socket.close();
    } else {
      this.error("this socket is not available");
    }
  };

  /**
   * 消息回調
   * @param {Object} message
   */
  $.onmessage = function (message) {};

  /**
   * 链接回调函数
   */
  $.onopen = function () {};

  /**
   * 关闭回调
   */
  $.onclose = function () {};

  /**
   * 异常回调
   */
  $.onerror = function () {};
})((ws = {}));
