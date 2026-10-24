// jforgame StructCodec's default (homogeneous collection) wire format.
// All numbers are big-endian. Field order comes exclusively from meta.schemas[].fields.
const MAX_BYTES = 1024 * 1024;
const utf8 = new TextEncoder();
const text = new TextDecoder('utf-8');
const numeric = {byte: ['Int8', 1], short: ['Int16', 2], int: ['Int32', 4],
  long: ['BigInt64', 8], float: ['Float32', 4], double: ['Float64', 8]};

class Writer {
  constructor() { this.buffer = new ArrayBuffer(256); this.offset = 0; }
  reserve(size) {
    const end = this.offset + size;
    if (end > MAX_BYTES) throw Error('Struct 消息体超过 1 MiB');
    if (end > this.buffer.byteLength) {
      const next = new ArrayBuffer(Math.max(end, Math.min(MAX_BYTES, this.buffer.byteLength * 2)));
      new Uint8Array(next).set(new Uint8Array(this.buffer)); this.buffer = next;
    }
  }
  number(name, value) {
    const [method, size] = numeric[name]; this.reserve(size);
    new DataView(this.buffer)[`set${method}`](this.offset, value); this.offset += size;
  }
  string(value) {
    const bytes = utf8.encode(value ?? ''); this.number('int', bytes.length); this.reserve(bytes.length);
    new Uint8Array(this.buffer, this.offset, bytes.length).set(bytes); this.offset += bytes.length;
  }
  finish() { return new Uint8Array(this.buffer.slice(0, this.offset)); }
}

class Reader {
  constructor(bytes) { this.bytes = bytes; this.view = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength); this.offset = 0; }
  require(size) {
    if (!Number.isInteger(size) || size < 0 || size > this.bytes.length - this.offset) throw Error('Struct 数据截断或长度非法');
  }
  number(name) {
    const [method, size] = numeric[name]; this.require(size);
    const value = this.view[`get${method}`](this.offset); this.offset += size; return value;
  }
  string() {
    const size = this.number('int'); this.require(size);
    const value = text.decode(this.bytes.subarray(this.offset, this.offset + size)); this.offset += size; return value;
  }
  count() {
    const size = this.number('short'); if (size < 0) throw Error('Struct 集合长度非法'); return size;
  }
}

function bytesOf(value) {
  if (value instanceof ArrayBuffer) return new Uint8Array(value);
  if (ArrayBuffer.isView(value)) return new Uint8Array(value.buffer, value.byteOffset, value.byteLength);
  throw Error('Struct 模式需要 WebSocket 二进制帧');
}

export class StructCodec {
  constructor(schemas) { this.schemas = schemas; }
  fields(ref) {
    const schema = this.schemas[ref];
    if (!schema) throw Error(`找不到 Struct 类型 ${ref}，请重启服务生成 meta`);
    return schema.fields;
  }
  encode(schema, body) {
    const writer = new Writer(); this.write(writer, {kind: 'object', ref: schema}, body, 0); return writer.finish();
  }
  decode(schema, data) {
    const reader = new Reader(bytesOf(data));
    const value = this.read(reader, {kind: 'object', ref: schema}, 0);
    if (reader.offset !== reader.bytes.length) throw Error('Struct 消息体有剩余字节，请检查 meta 与服务端字段顺序');
    return value;
  }
  write(out, type, value, depth) {
    if (depth > 64) throw Error('Struct 嵌套超过 64 层');
    if (type.kind === 'object') {
      if (value == null || typeof value !== 'object' || Array.isArray(value)) throw Error(`Struct 对象 ${type.ref} 不能为空`);
      for (const field of this.fields(type.ref)) {
        const item = Object.hasOwn(value, field.name) && value[field.name] !== undefined
          ? value[field.name] : field.hasDefault ? field.defaultValue : undefined;
        try { this.write(out, field.type, item, depth + 1); }
        catch (error) { throw Error(`${field.name}：${error.message}`); }
      }
    } else if (type.kind === 'string') {
      if (value != null && typeof value !== 'string') throw Error('请输入字符串');
      out.string(value);
    } else if (type.kind === 'list' || type.kind === 'map') {
      if (type.kind === 'list' && value != null && !Array.isArray(value)) throw Error('请输入数组');
      if (type.kind === 'map' && value != null && (typeof value !== 'object' || Array.isArray(value))) throw Error('请输入 Map');
      // Struct's MapCodec always serializes keys as UTF-8 strings, regardless of Java generic key type.
      const entries = value == null ? [] : type.kind === 'list' ? value : Object.entries(value);
      if (entries.length > 32767) throw Error('Struct 集合最多 32767 个元素');
      out.number('short', entries.length);
      for (const entry of entries) {
        const item = type.kind === 'map' ? entry[1] : entry;
        if (item == null) throw Error('Struct 集合元素不能为 null');
        if (type.kind === 'map') out.string(entry[0]);
        this.write(out, type.element, item, depth + 1);
      }
    } else if (type.kind === 'boolean') {
      if (typeof value !== 'boolean') throw Error('Struct boolean 不能为空且必须为布尔值');
      out.number('byte', value ? 1 : 0);
    } else if (['integer', 'long', 'number'].includes(type.kind) && numeric[type.name]) {
      if (value == null) throw Error(`Struct ${type.name} 不能为空`);
      if (type.kind === 'long') {
        if (typeof value === 'number' && !Number.isSafeInteger(value)) throw Error('long 请使用整数文本或 BigInt');
        if (!/^-?\d+$/.test(String(value))) throw Error('请输入 long 整数');
        value = BigInt(value);
        if (value < -9223372036854775808n || value > 9223372036854775807n) throw Error('超出 long 范围');
      } else {
        value = Number(value);
        if (!Number.isFinite(value)) throw Error('请输入有限数值');
        const bits = numeric[type.name][1] * 8;
        if (type.kind === 'integer' && (!Number.isInteger(value) || value < -(2 ** (bits - 1)) || value >= 2 ** (bits - 1))) throw Error('整数超出范围');
        if (type.name === 'float' && !Number.isFinite(Math.fround(value))) throw Error('超出 float 范围');
      }
      out.number(type.name, value);
    } else throw Error(`StructCodec 不支持 ${type.name ?? type.kind} 类型`);
  }
  read(input, type, depth) {
    if (depth > 64) throw Error('Struct 嵌套超过 64 层');
    if (type.kind === 'object') {
      const value = Object.create(null);
      for (const field of this.fields(type.ref)) value[field.name] = this.read(input, field.type, depth + 1);
      return value;
    }
    if (type.kind === 'string') return input.string();
    if (type.kind === 'boolean') return input.number('byte') > 0;
    if (['integer', 'long', 'number'].includes(type.kind) && numeric[type.name]) {
      const value = input.number(type.name); return typeof value === 'bigint' ? value.toString() : value;
    }
    if (type.kind === 'list' || type.kind === 'map') {
      const size = input.count(); const value = type.kind === 'list' ? [] : Object.create(null);
      for (let i = 0; i < size; i++) {
        const key = type.kind === 'map' ? input.string() : i;
        value[key] = this.read(input, type.element, depth + 1);
      }
      return value;
    }
    throw Error(`StructCodec 不支持 ${type.name ?? type.kind} 类型`);
  }
}

export function encodeStructFrame(codec, protocol, index, body) {
  const payload = codec.encode(protocol.schema, body);
  const frame = new Uint8Array(12 + payload.length); const header = new DataView(frame.buffer);
  header.setInt32(0, frame.length); header.setInt32(4, index); header.setInt32(8, protocol.cmd);
  frame.set(payload, 12); return frame.buffer;
}

export function decodeStructFrame(codec, byCmd, data) {
  const bytes = bytesOf(data);
  if (bytes.length < 12) throw Error('Struct 帧不足 12 字节');
  const header = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength);
  if (header.getInt32(0) !== bytes.length) throw Error('Struct 帧长度与消息头不符');
  const frame = {index: header.getInt32(4), cmd: header.getInt32(8)};
  const protocol = byCmd.get(frame.cmd);
  if (!protocol) throw Error(`未知 Struct cmd ${frame.cmd}，请刷新 meta`);
  return {frame, body: codec.decode(protocol.schema, bytes.subarray(12))};
}
