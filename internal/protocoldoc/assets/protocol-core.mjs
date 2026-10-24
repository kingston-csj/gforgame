import {StructCodec, encodeStructFrame, decodeStructFrame} from './struct-codec.mjs';

const utf8 = new TextEncoder();
const text = new TextDecoder('utf-8', {fatal: true});

// WebSocket binary framing is independent of the message body's JSON/Struct codec.
// Both use the server's big-endian header: total length, request index, command.
export function encodeJsonFrame(protocol, index, body) {
  const payload = utf8.encode(stringifyJson(body));
  const bytes = new Uint8Array(12 + payload.length);
  const header = new DataView(bytes.buffer);
  header.setInt32(0, bytes.length);
  header.setInt32(4, index);
  header.setInt32(8, protocol.cmd);
  bytes.set(payload, 12);
  return bytes.buffer;
}

export function decodeJsonFrame(data) {
  // Accept legacy text envelopes when connected to an older text-frame endpoint.
  if (typeof data === 'string') {
    const frame = parseJson(data);
    return {frame, body: typeof frame.msg === 'string' ? parseJson(frame.msg) : frame.msg};
  }
  const bytes = data instanceof ArrayBuffer ? new Uint8Array(data)
    : ArrayBuffer.isView(data) ? new Uint8Array(data.buffer, data.byteOffset, data.byteLength) : null;
  if (!bytes) throw Error('JSON 模式需要 WebSocket 二进制帧或文本帧');
  if (bytes.length < 12) throw Error('JSON 二进制帧不足 12 字节');
  const header = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength);
  if (header.getInt32(0) !== bytes.length) throw Error('JSON 二进制帧长度与消息头不符');
  return {frame: {index: header.getInt32(4), cmd: header.getInt32(8)},
    body: parseJson(text.decode(bytes.subarray(12)))};
}

// Preserve integer values outside JavaScript's safe range in both directions.
export function parseJson(text) {
  return JSON.parse(text.replace(/"(?:\\.|[^"\\])*"|-?\d+(?:\.\d+)?(?:e[+-]?\d+)?/gi, token => {
    if (/^-?\d+$/.test(token) && !Number.isSafeInteger(Number(token))) return JSON.stringify(token);
    return token;
  }));
}

export function stringifyJson(value) {
  if (typeof value === 'bigint') return value.toString();
  if (Array.isArray(value)) return `[${value.map(item => stringifyJson(item) ?? 'null').join(',')}]`;
  if (value !== null && typeof value === 'object') return `{${Object.entries(value)
    .filter(([, item]) => item !== undefined)
    .map(([key, item]) => `${JSON.stringify(key)}:${stringifyJson(item)}`).join(',')}}`;
  return JSON.stringify(value);
}

export function scalarValue(type, value) {
  if (type.kind === 'string') return value;
  if (type.kind === 'enum') {
    if (!type.values.includes(value)) throw Error('请选择有效枚举值');
    return value;
  }
  if (type.kind === 'boolean') return value === 'true';
  if (type.kind === 'long') {
    if (!/^-?\d+$/.test(value)) throw Error('请输入 long 整数');
    const number = BigInt(value);
    if (number < -9223372036854775808n || number > 9223372036854775807n) throw Error('超出 long 范围');
    return number;
  }
  if (type.kind === 'integer') {
    if (!/^-?\d+$/.test(value)) throw Error('请输入整数');
    const number = Number(value);
    const [min, max] = {byte: [-128, 127], short: [-32768, 32767]}[type.name] ?? [-2147483648, 2147483647];
    if (!Number.isSafeInteger(number) || number < min || number > max) throw Error('整数超出范围');
    return number;
  }
  if (type.kind === 'number') {
    if (!value.trim() || !Number.isFinite(Number(value))) throw Error('请输入有限数值');
    const number = Number(value);
    if (type.name === 'float' && Math.abs(number) > 3.4028234663852886e38) throw Error('超出 float 范围');
    return number;
  }
  throw Error(`暂不支持输入 ${type.name} 类型`);
}

export function matchingRequests(protocols, query, module) {
  const term = query.trim().toLowerCase();
  return protocols.filter(p => p.name.startsWith('Req') && p.testable &&
    (!module || p.module === module) && p.name.toLowerCase().includes(term));
}

// Match actual handler return types first; conventional names cover requests without handlers.
export function orderedProtocols(protocols) {
  const sorted = [...protocols].sort((a, b) => (a.cmd ?? Infinity) - (b.cmd ?? Infinity) || a.name.localeCompare(b.name));
  const byId = new Map(protocols.map(p => [p.id, p]));
  const emitted = new Set();
  const result = [];
  const append = protocol => {
    if (protocol && !emitted.has(protocol.id)) { emitted.add(protocol.id); result.push(protocol); }
  };
  for (const req of sorted.filter(p => p.direction === 'REQ')) {
    append(req);
    const conventionalId = req.id.slice(0, req.id.lastIndexOf('.') + 1) + 'Res' + req.name.slice(3);
    const response = byId.get(req.response) || byId.get(conventionalId);
    if (response?.direction === 'RES') append(response);
  }
  for (const direction of ['RES', 'PUSH', 'VO']) {
    for (const protocol of sorted.filter(p => p.direction === direction)) append(protocol);
  }
  return result;
}

export class ProtocolSocket {
  constructor(onMessage, onClose, socketFactory = url => new WebSocket(url), timeout = 15000) {
    this.onMessage = onMessage;
    this.onClose = onClose;
    this.socketFactory = socketFactory;
    this.timeout = timeout;
    this.pending = new Map();
    this.index = 0;
    this.socket = null;
    this.codec = 'json';
    this.byCmd = new Map();
  }
  configure(meta) {
    this.struct = new StructCodec(meta.schemas);
    this.byCmd = new Map(meta.modules.flatMap(module => module.protocols)
      .filter(protocol => protocol.cmd != null).map(protocol => [protocol.cmd, protocol]));
  }
  setCodec(codec) {
    if (!['json', 'struct'].includes(codec)) throw Error('未知编解码类型');
    if (this.codec === codec) return;
    this.close(); this.codec = codec;
  }
  connect(url) {
    if (this.socket) throw Error('连接已存在，请先断开');
    return new Promise((resolve, reject) => {
      const socket = this.socketFactory(url);
      this.socket = socket;
      socket.binaryType = 'blob';
      const timer = setTimeout(() => { reject(Error('连接超时')); this.close(); }, this.timeout);
      this.cancelConnect = () => { clearTimeout(timer); reject(Error('连接已断开')); };
      socket.onopen = () => {
        if (this.socket !== socket) return;
        clearTimeout(timer); this.cancelConnect = null; resolve();
      };
      socket.onerror = () => {
        if (this.socket !== socket) return;
        clearTimeout(timer); reject(Error('WebSocket 连接失败，请检查服务地址')); this.close();
      };
      socket.onclose = () => {
        clearTimeout(timer);
        reject(Error('连接已断开'));
        if (this.socket === socket) {
          this.socket = null;
          this.cancelConnect = null;
          this.rejectPending();
          this.onClose();
        }
      };
      socket.onmessage = async event => {
        if (this.socket !== socket) return;
        try {
          let frame, body;
          const binary = event.data instanceof Blob || event.data instanceof ArrayBuffer || ArrayBuffer.isView(event.data);
          if (binary) {
            const data = event.data instanceof Blob ? await event.data.arrayBuffer() : event.data;
            // Blob conversion may finish after disconnecting or switching codecs.
            if (this.socket !== socket) return;
            if (this.codec === 'struct') {
              ({frame, body} = decodeStructFrame(this.struct, this.byCmd, data));
            } else {
              ({frame, body} = decodeJsonFrame(data));
            }
          } else {
            // Text frames keep the original JSON envelope, regardless of body codec.
            if (typeof event.data !== 'string') throw Error('未知的 WebSocket 数据帧类型');
            ({frame, body} = decodeJsonFrame(event.data));
          }
          const candidate = this.pending.get(frame.index);
          const pending = candidate && (candidate.responseCmd == null || candidate.responseCmd === frame.cmd) &&
            this.byCmd.get(frame.cmd)?.direction !== 'PUSH' ? candidate : null;
          const elapsed = pending ? Date.now() - pending.started : null;
          this.onMessage(frame, body, elapsed);
          if (pending) {
            clearTimeout(pending.timer);
            this.pending.delete(frame.index);
            pending.resolve({frame, body, elapsed});
          }
        } catch (error) {
          if (this.socket !== socket) return;
          this.onMessage({cmd: 0, index: 0, error: true}, {message: error.message}, null);
        }
      };
    });
  }
  request(protocol, body) {
    if (this.socket?.readyState !== 1) return Promise.reject(Error('连接未就绪'));
    if (!protocol.testable || !protocol.name.startsWith('Req') || !Number.isInteger(protocol.cmd)) {
      return Promise.reject(Error('请选择可测试的 Req 协议'));
    }
    const index = this.index = this.index >= 2147483646 ? 1 : this.index + 1;
    return new Promise((resolve, reject) => {
      const timer = setTimeout(() => {
        this.pending.delete(index);
        reject(Error(`${protocol.name} 请求超时（15 秒），请查看推送或服务器日志`));
      }, this.timeout);
      const response = [...this.byCmd.values()].find(item => item.id === protocol.response ||
        item.name === 'Res' + protocol.name.slice(3));
      this.pending.set(index, {resolve, reject, timer, started: Date.now(), responseCmd: response?.cmd});
      try {
        this.socket.send(this.codec === 'struct'
          ? encodeStructFrame(this.struct, protocol, index, body)
          : encodeJsonFrame(protocol, index, body));
        this.onMessage({cmd: protocol.cmd, index, sent: true}, body, null);
      } catch (error) {
        clearTimeout(timer);
        this.pending.delete(index);
        reject(error);
      }
    });
  }
  rejectPending() {
    for (const entry of this.pending.values()) {
      clearTimeout(entry.timer);
      entry.reject(Error('连接已断开'));
    }
    this.pending.clear();
  }
  close() {
    this.cancelConnect?.(); this.cancelConnect = null;
    const socket = this.socket;
    this.socket = null;
    this.rejectPending();
    if (socket) {
      socket.close();
      this.onClose();
    }
  }
}
