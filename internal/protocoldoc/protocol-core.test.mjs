import {test} from 'node:test';
import assert from 'node:assert/strict';
import {ProtocolSocket, parseJson, stringifyJson, encodeJsonFrame, decodeJsonFrame} from './assets/protocol-core.mjs';

test('64 位整数收发保持精度', () => {
  assert.equal(parseJson('{"id":9223372036854775807}').id, '9223372036854775807');
  assert.equal(stringifyJson({id: 9223372036854775807n}), '{"id":9223372036854775807}');
});

const req = {id: 'ReqPlayerLogin', name: 'ReqPlayerLogin', cmd: 103, testable: true, response: 'ResPlayerLogin', direction: 'REQ'};
const res = {id: 'ResPlayerLogin', name: 'ResPlayerLogin', cmd: 154, direction: 'RES'};

async function connected(timeout = 500) {
  const records = [];
  const fake = {readyState: 1, send(data) { this.sent = data; }, close() {}};
  const socket = new ProtocolSocket((...args) => records.push(args), () => {}, () => fake, timeout);
  socket.configure({schemas: {}, modules: [{protocols: [req, res, {cmd: 155, direction: 'PUSH'}]}]});
  const ready = socket.connect('ws://localhost/ws'); fake.onopen(); await ready;
  return {socket, fake, records};
}

test('Go JSON 登录帧、响应配对和登录推送', async () => {
  const {socket, fake, records} = await connected();
  try {
    const pending = socket.request(req, {PlayerId: 'tester', serverId: 14});
    assert.ok(fake.sent instanceof ArrayBuffer);
    assert.deepEqual(decodeJsonFrame(fake.sent), {frame: {cmd: 103, index: 1}, body: {PlayerId: 'tester', serverId: 14}});
    fake.onmessage({data: '{"cmd":155,"index":1,"msg":{}}'});
    assert.equal(socket.pending.size, 1);
    await fake.onmessage({data: new Blob([encodeJsonFrame(res, 1, {code: 0, playerId: 'tester'})])});
    assert.equal((await pending).body.playerId, 'tester');
    assert.equal(records.length, 3);
    assert.equal(socket.pending.size, 0);
  } finally { socket.close(); }
});

test('断开连接拒绝待处理请求', async () => {
  const {socket} = await connected();
  const pending = socket.request(req, {});
  socket.close();
  await assert.rejects(pending, /连接已断开/);
  assert.equal(socket.pending.size, 0);
});

test('请求超时不会重发', async () => {
  const {socket} = await connected(10);
  try { await assert.rejects(socket.request(req, {}), /请求超时/); }
  finally { socket.close(); }
});
