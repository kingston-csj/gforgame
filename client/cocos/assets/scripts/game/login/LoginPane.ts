import { _decorator, Component, Node, EditBox, Button, director } from 'cc';
import { WebSocketClient } from '../../net/WebSocketClient';
import GameContext from '../../GameContext';
import RespLogin from '../../net/RespLogin';
import ReqLogin from '../../net/ReqLogin';
import { MessageDispatch } from '../../MessageDispatch';
const { ccclass, property } = _decorator;

@ccclass('LoginScene')
export class LoginScene extends Component {
  @property(EditBox)
  usernameInput: EditBox = null!;

  @property(EditBox)
  passwordInput: EditBox = null!;

  @property(Button)
  loginButton: Button = null!;

  @property(Button)
  logoutButton: Button = null!;

  start() {
    // 注册登录按钮点击事件
    this.passwordInput.inputFlag = EditBox.InputFlag.PASSWORD;
    this.loginButton.node.on(Button.EventType.CLICK, this.onLoginClick, this);

    MessageDispatch.register(RespLogin.cmd, (msg: RespLogin) => {
      console.log('登录成功', msg);
    });

    GameContext.instance.WebSocketClient.connect('ws://127.0.0.1:9527/ws');
  }

  onLoginClick() {
    const username = this.usernameInput.string;
    const password = this.passwordInput.string;

    // 这里添加登录验证逻辑
    if (username && password) {
      console.log('登录信息:', { username, password });
      // TODO: 实现实际的登录逻辑
      // 登录成功后跳转到游戏主场景
      // director.loadScene('GameMain');

      GameContext.instance.WebSocketClient.sendMessage(
        ReqLogin.cmd,
        {
          id: username,
          pwd: password,
        },
        (msg: RespLogin) => {
          console.log('登录成功');
        }
      );
    } else {
      console.log('请输入用户名和密码');
    }
  }
}
