import {parseJson, scalarValue, matchingRequests, orderedProtocols, ProtocolSocket} from './protocol-core.mjs';

const $ = id => document.getElementById(id);
const el = (tag, text, className) => {
  const node = document.createElement(tag);
  if (text !== undefined) node.textContent = text;
  if (className) node.className = className;
  return node;
};
let meta, protocols = [], moduleId = '', selected, form, online = false, loggingIn = false;
let sessionVersion = 0;
const byCmd = new Map();
const socket = new ProtocolSocket(receive, () => {
  sessionVersion++;
  online = false; loggingIn = false; updateConnection('未连接');
  $('send-status').textContent = '';
});
function showError(error) {
  $('error').textContent = error?.message || String(error);
  $('error').hidden = false;
}
function clearError() { $('error').hidden = true; }
function tab(name) {
  for (const id of ['docs', 'test']) {
    $(`${id}-panel`).hidden = id !== name;
    $(`${id}-tab`).setAttribute('aria-selected', String(id === name));
  }
}
$('docs-tab').onclick = () => tab('docs');
$('test-tab').onclick = () => tab('test');

function renderModules() {
  $('modules').replaceChildren();
  const items = [{id: '', name: '全部模块', protocols}, ...meta.modules];
  for (const item of items) {
    const button = el('button', undefined, item.id === moduleId ? 'active' : '');
    button.append(el('span', item.id ? `${item.name} · ${item.id}` : item.name), el('small', item.protocols.length));
    button.onclick = () => {
      moduleId = item.id;
      renderModules(); renderDocs(); renderOptions();
      if (selected && moduleId && selected.module !== moduleId && !$('all-modules').checked) clearSelection();
    };
    $('modules').append(button);
  }
}
function typeView(type, visited = new Set()) {
  const node = el('div', undefined, 'type');
  if (type.ref && meta.schemas[type.ref]) {
    const button = el('button', type.name, 'type-ref');
    const content = el('div', undefined, 'nested-schema'); content.hidden = true;
    button.onclick = () => {
      if (!content.hasChildNodes()) {
        if (visited.has(type.ref)) content.append(el('span', '循环引用：' + type.name));
        else content.append(fieldTable(meta.schemas[type.ref].fields, new Set([...visited, type.ref])));
      }
      content.hidden = !content.hidden;
    };
    node.append(button, content);
  } else {
    node.append(el('span', type.name));
    if (type.element) node.append(typeView(type.element, visited));
    if (type.kind === 'enum') node.append(el('div', type.values.join(' / '), 'hint'));
  }
  return node;
}
function fieldTable(fields, visited = new Set(), editable = false) {
  const table = el('table', undefined, 'field-table' + (editable ? ' editable' : ''));
  const header = el('tr');
  for (const text of ['字段名', '类型', '说明', ...(editable ? ['参数值'] : [])]) header.append(el('th', text));
  const head = el('thead'); head.append(header); table.append(head);
  const body = el('tbody'); table.append(body);
  for (const field of fields) {
    const row = el('tr');
    const typeCell = el('td'); typeCell.append(typeView(field.type, visited));
    if (field.nullable) typeCell.append(el('small', '可空', 'hint'));
    row.append(el('td', field.name), typeCell, el('td', field.desc || '—')); body.append(row);
  }
  if (!fields.length) {
    const row = el('tr'); const cell = el('td', '无字段', 'hint');
    cell.colSpan = editable ? 4 : 3; row.append(cell); body.append(row);
  }
  return table;
}
function renderDocs() {
  const query = $('doc-search').value.trim().toLowerCase();
  const group = meta.modules.find(item => item.id === moduleId);
  $('module-title').textContent = group ? `${group.name} · ${group.id}` : '全部协议';
  const target = $('documents'); target.replaceChildren();
  let count = 0;
  for (const module of meta.modules.filter(item => !moduleId || item.id === moduleId)) {
    const matches = orderedProtocols(module.protocols).filter(p => `${p.name} ${p.cmd ?? ''} ${p.desc}`.toLowerCase().includes(query));
    if (!matches.length) continue;
    if (!moduleId) target.append(el('h2', `${module.name} · ${module.id}`, 'module-heading'));
    for (const protocol of matches) {
      count++;
      const details = el('details', undefined, 'protocol'); details.id = protocol.id;
      const summary = el('summary');
      summary.append(el('span', protocol.direction, `badge ${protocol.direction}`), el('strong', protocol.name), el('span', protocol.desc, 'description'));
      const content = el('div', undefined, 'content');
      if (protocol.direction !== 'VO') content.append(el('p', `cmd：${protocol.cmd ?? '未注册'}`, 'cmd'));
      content.append(fieldTable(meta.schemas[protocol.schema]?.fields || [], new Set([protocol.schema])));
      if (protocol.testable) {
        const test = el('button', '在线测试');
        test.onclick = () => {
          tab('test'); $('all-modules').checked = true; $('request-search').value = protocol.name;
          selectRequest(protocol); renderOptions();
        };
        content.append(test);
      }
      details.append(summary, content); target.append(details);
      details.ontoggle = () => { if (details.open) history.replaceState(null, '', '#' + encodeURIComponent(protocol.id)); };
    }
  }
  if (!count) target.append(el('p', '没有匹配的协议', 'empty'));
}
$('doc-search').oninput = () => meta && renderDocs();

// Explicit omission/null/value modes distinguish defaults from empty inputs.
function fieldEditor(field, initial, supplied = false, depth = 0) {
  const node = el('div', undefined, 'value-editor');
  const mode = el('select', undefined, 'mode');
  const hasInitial = supplied || field.hasDefault;
  const value = supplied ? initial : field.defaultValue;
  for (const [id, text] of [['value', '填写参数'], ['omit', '不传此字段'], ...(field.nullable ? [['null', 'null']] : [])]) {
    const option = el('option', text); option.value = id; mode.append(option);
  }
  mode.value = hasInitial ? (value == null && field.nullable ? 'null' : 'value') : 'omit';
  const control = valueEditor(field.type, value, depth);
  const update = () => { control.node.hidden = mode.value !== 'value'; };
  mode.onchange = update; mode.setAttribute('aria-label', `${field.name} 参数模式`);
  node.append(mode, control.node); update();
  return {node, read() {
    if (mode.value === 'omit') return undefined;
    if (mode.value === 'null') return null;
    try { return control.read(); } catch (error) { throw Error(`${field.name}：${error.message}`); }
  }};
}
function valueEditor(type, initial, depth) {
  if (['string', 'integer', 'long', 'number', 'boolean', 'enum'].includes(type.kind)) {
    const node = el(type.kind === 'boolean' || type.kind === 'enum' ? 'select' : 'input', undefined, 'scalar');
    if (node.tagName === 'SELECT') {
      for (const value of type.kind === 'boolean' ? ['false', 'true'] : type.values) {
        const option = el('option', value); option.value = value; node.append(option);
      }
    } else {
      node.type = 'text';
      if (type.kind !== 'string') node.inputMode = type.kind === 'number' ? 'decimal' : 'numeric';
      node.placeholder = type.name;
    }
    if (initial !== null && initial !== undefined) node.value = String(initial);
    return {node, read: () => scalarValue(type, node.value)};
  }
  if (type.kind === 'object') {
    const node = el('details', undefined, 'object-editor'); node.append(el('summary', type.name));
    let editor;
    const ensure = () => {
      if (editor) return;
      if (depth >= 12) throw Error('嵌套超过 12 层，请检查循环引用');
      const schema = meta.schemas[type.ref];
      if (!schema) throw Error(`找不到类型 ${type.name}`);
      editor = objectForm(schema, initial, depth + 1); node.append(editor.node);
    };
    node.ontoggle = () => { if (node.open) { try { ensure(); } catch (error) { showError(error); node.open = false; } } };
    return {node, read: () => { ensure(); return editor.read(); }};
  }
  if (type.kind === 'list' || type.kind === 'map') {
    const node = el('div', undefined, 'collection-editor');
    const rows = el('div'); node.append(rows);
    const entries = [];
    const add = (value, key) => {
      const row = el('div', undefined, 'collection-row');
      const item = fieldEditor({name: '元素', type: type.element, nullable: !['integer', 'number', 'long', 'boolean'].includes(type.element.kind), hasDefault: false}, value, value !== undefined, depth + 1);
      const keyInput = el('input'); keyInput.placeholder = `键：${type.key?.name ?? 'string'}`; keyInput.value = key ?? '';
      if (type.kind === 'map') row.append(keyInput);
      const remove = el('button', '×'); remove.title = '删除元素';
      const entry = {item, keyInput}; entries.push(entry);
      remove.onclick = () => { entries.splice(entries.indexOf(entry), 1); row.remove(); };
      row.append(item.node, remove); rows.append(row);
    };
    const addButton = el('button', type.kind === 'map' ? '+ 添加键值' : '+ 添加元素');
    addButton.onclick = () => add(); node.append(addButton);
    if (type.kind === 'list' && Array.isArray(initial)) initial.forEach(value => add(value));
    if (type.kind === 'map' && initial && typeof initial === 'object') Object.entries(initial).forEach(([key, value]) => add(value, key));
    return {node, read() {
      if (depth >= 12) throw Error('嵌套超过 12 层');
      if (type.kind === 'list') return entries.map(({item}, index) => {
        const value = item.read();
        if (value === undefined) throw Error(`第 ${index + 1} 个元素尚未填写`);
        return value;
      });
      const result = Object.create(null);
      for (const {item, keyInput} of entries) {
        if (!keyInput.value) throw Error('Map 键不能为空');
        const key = String(scalarValue(type.key, keyInput.value));
        if (Object.hasOwn(result, key)) throw Error(`Map 键重复：${key}`);
        const value = item.read();
        if (value === undefined) throw Error(`Map 值尚未填写：${key}`);
        result[key] = value;
      }
      return result;
    }};
  }
  return {node: el('p', `暂不支持 ${type.name} 类型`, 'hint'), read() { throw Error(`暂不支持 ${type.name} 类型`); }};
}
function objectForm(schema, initial, depth = 0) {
  const fields = schema.fields.filter(field => socket.codec === 'struct' || field.json !== false);
  const node = fieldTable(fields, new Set(), true);
  const readers = [];
  [...node.tBodies[0].rows].forEach((row, index) => {
    const field = fields[index]; if (!field) return;
    const supplied = initial != null && Object.hasOwn(initial, field.name);
    const control = fieldEditor(field, supplied ? initial[field.name] : undefined, supplied, depth);
    const cell = el('td'); cell.append(control.node); row.append(cell); readers.push([field.name, control]);
  });
  return {node, read() {
    const value = Object.create(null);
    for (const [name, reader] of readers) {
      const item = reader.read(); if (item !== undefined) value[name] = item;
    }
    return value;
  }};
}
function clearSelection() {
  selected = null; form = null;
  $('request-detail').replaceChildren(el('p', '选择协议后填写参数。', 'hint'));
  $('reset').disabled = true; $('send').disabled = true;
}
function selectRequest(protocol) {
  selected = protocol; form = objectForm(meta.schemas[protocol.schema]);
  $('request-detail').replaceChildren(el('h2', protocol.name), el('p', protocol.desc || '—', 'hint'), el('p', `cmd：${protocol.cmd}`, 'cmd'), form.node);
  $('reset').disabled = false; $('send').disabled = !online; $('send-status').textContent = '';
}
function renderOptions() {
  const options = matchingRequests(protocols, $('request-search').value, $('all-modules').checked ? '' : moduleId);
  $('request-options').replaceChildren();
  for (const protocol of options) {
    const button = el('button', protocol.name, selected?.id === protocol.id ? 'selected' : '');
    button.append(el('small', `${protocol.module} · cmd ${protocol.cmd} · ${protocol.desc}`));
    button.onclick = () => { selectRequest(protocol); renderOptions(); }; $('request-options').append(button);
  }
  if (!options.length) $('request-options').append(el('p', '没有匹配的 Req 协议', 'hint'));
}
$('request-search').oninput = () => meta && renderOptions();
$('all-modules').onchange = () => {
  if (selected && moduleId && selected.module !== moduleId && !$('all-modules').checked) clearSelection();
  if (meta) renderOptions();
};
$('reset').onclick = () => selected && selectRequest(selected);
function endpoint() {
  if (meta.connection.url) return meta.connection.url;
  const url = new URL(location.href);
  url.protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'; url.port = String(meta.connection.port);
  url.pathname = meta.connection.path; url.search = ''; url.hash = ''; return url.href;
}
function updateConnection(text) {
  $('connection-status').textContent = text;
  $('connection-status').className = 'status' + (online ? ' online' : loggingIn ? ' busy' : '');
  $('login').disabled = !meta?.connection.testEnabled || online || loggingIn;
  $('player-id').disabled = online || loggingIn;
  $('disconnect').disabled = !online && !loggingIn; $('send').disabled = !online || !selected;
}
function loginProtocol(name) {
  const protocol = protocols.find(item => item.name === name && item.testable);
  if (!protocol) throw Error(`缺少登录协议 ${name}`); return protocol;
}
function checkLogin(result, name) {
  if (!result.body || Number(result.body.code ?? 0) !== 0) throw Error(`${name}失败，错误码：${result.body?.code ?? '未知'}`);
}
$('login').onclick = async () => {
  clearError(); const playerId = $('player-id').value.trim();
  if (!playerId) { showError(Error('请输入 playerId')); return; }
  loggingIn = true; updateConnection('正在连接…');
  const version = ++sessionVersion;
  try {
    const playerProtocol = loginProtocol('ReqPlayerLogin');
    await socket.connect(endpoint());
    if (version !== sessionVersion) return;
    updateConnection('正在登录角色…');
    const fields = meta.schemas[playerProtocol.schema].fields;
    const playerField = fields.find(field => field.name.toLowerCase() === 'playerid')?.name ?? 'playerId';
    const result = await socket.request(playerProtocol, {[playerField]: playerId, serverId: meta.connection.serverId});
    if (version !== sessionVersion) return;
    checkLogin(result, '角色登录');
    if (!socket.socket || !loggingIn) throw Error('登录连接已断开');
    online = true; loggingIn = false; updateConnection(`已登录 · ${playerId}`);
  } catch (error) {
    if (version !== sessionVersion) return;
    socket.close(); online = false; loggingIn = false; updateConnection('未连接'); showError(error);
  }
};
$('disconnect').onclick = () => socket.close();
$('codec').onchange = () => {
  socket.setCodec($('codec').value);
  try { localStorage.setItem('protocol-doc-codec', socket.codec); } catch { /* Storage may be disabled. */ }
  clearError(); updateConnection('未连接');
  if (selected) selectRequest(selected);
};
$('player-id').onkeydown = event => { if (event.key === 'Enter' && !$('login').disabled) $('login').click(); };
$('send').onclick = async () => {
  clearError(); if (!online || !selected || !form) return;
  $('send').disabled = true; const current = selected;
  const version = sessionVersion;
  try {
    const body = form.read(); $('send-status').textContent = '等待响应…';
    const result = await socket.request(current, body);
    if (version !== sessionVersion) return;
    $('send-status').textContent = `${current.name} 已收到响应 · ${result.elapsed} ms`;
    if (result.body?.code != null && Number(result.body.code) !== 0) showError(Error(`请求失败，错误码：${result.body.code}`));
  } catch (error) {
    if (version !== sessionVersion) return;
    $('send-status').textContent = '发送未完成'; showError(error); appendLog('ERROR', current.name, {message: error.message});
  } finally { if (version === sessionVersion) $('send').disabled = !online || !selected; }
};
function dataTree(value, name = '') {
  if (value !== null && typeof value === 'object') {
    const details = el('details', undefined, 'data-tree'); const entries = Object.entries(value);
    details.append(el('summary', `${name || '字段'} ${Array.isArray(value) ? `[${entries.length}]` : `{${entries.length}}`}`));
    details.ontoggle = () => {
      if (details.open && details.childElementCount === 1) for (const [key, item] of entries) details.append(dataTree(item, key));
    };
    return details;
  }
  return el('div', `${name ? name + '：' : ''}${value === null ? 'null' : String(value)}`, 'data-leaf');
}
function appendLog(kind, name, body, extra = '') {
  const record = el('details', undefined, 'message'); const summary = el('summary');
  summary.append(el('time', new Date().toLocaleTimeString()), el('span', `${kind} ${name} ${extra}`, kind));
  record.append(summary, dataTree(body)); $('messages').append(record);
  while ($('messages').childElementCount > 200) $('messages').firstElementChild.remove();
  if (!$('pause-scroll').checked) $('messages').scrollTop = $('messages').scrollHeight;
}
function receive(frame, body, elapsed) {
  const protocol = byCmd.get(frame.cmd);
  const kind = frame.error ? 'ERROR' : frame.sent ? 'REQ' : elapsed !== null ? 'RES' : 'PUSH';
  appendLog(kind, protocol?.name || `cmd ${frame.cmd}`, body, `#${frame.index}${elapsed !== null ? ` · ${elapsed} ms` : ''}`);
  if (!frame.sent && protocol?.name === 'PushPlayerForceLogout') {
    socket.close(); showError(Error('当前玩家已被强制下线，请重新登录'));
  }
}
$('clear-log').onclick = () => $('messages').replaceChildren();
window.addEventListener('beforeunload', () => socket.close());
async function init() {
  try {
    const response = await fetch('./meta.json', {cache: 'no-store'});
    if (!response.ok) throw Error(`协议 meta 加载失败（HTTP ${response.status}），请确认 game.protocol-doc 已设为 true 且服务启动完成`);
    meta = parseJson(await response.text()); protocols = meta.modules.flatMap(module => module.protocols);
    socket.configure(meta);
    try {
      const allowed = meta.connection.codecs ?? ['json'];
      for (const option of [...$('codec').options]) if (!allowed.includes(option.value)) option.remove();
      $('codec').disabled = allowed.length === 1;
      const savedCodec = allowed.includes(localStorage.getItem('protocol-doc-codec')) ? localStorage.getItem('protocol-doc-codec') : 'json';
      if (savedCodec === 'struct') { socket.setCodec(savedCodec); $('codec').value = savedCodec; }
    } catch { /* Storage may be disabled. */ }
    for (const protocol of protocols) if (protocol.cmd != null) byCmd.set(protocol.cmd, protocol);
    $('meta-status').textContent = `${meta.modules.length} 个模块 · ${protocols.length} 个协议`;
    $('endpoint').textContent = meta.connection.testEnabled ? `连接地址：${endpoint()}` : '当前为网关模式，请配置 game.protocol-doc-options.websocket-url 为客户端网关入口地址';
    renderModules(); renderDocs(); renderOptions(); updateConnection('未连接');
    if (location.hash) {
      const target = document.getElementById(decodeURIComponent(location.hash.slice(1)));
      if (target) { target.open = true; target.scrollIntoView({block: 'center'}); }
    }
  } catch (error) { $('meta-status').textContent = '加载失败'; showError(error); }
}
init();
