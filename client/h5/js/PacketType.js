/**
 * 与服务端的通信协议绑定
 */

PacketType = {}

PacketType.ReqAccountLogin = "2001";

PacketType.ResAccountLogin = "2002";

var self = PacketType;

var msgHandler = {}

PacketType.bind = function(msgId, handler) {
	msgHandler[msgId] = handler
}

self.bind(self.ResAccountLogin, function(resp) {
	alert("角色登录成功-->" + resp)
})

PacketType.handle = function(msgId, msg) {
	msgHandler[msgId](msg);
}
