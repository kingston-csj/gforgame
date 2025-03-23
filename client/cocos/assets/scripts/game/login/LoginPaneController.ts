import { _decorator, Component, Node, EditBox, Button, director } from 'cc';
import GameContext from '../../GameContext';
import RespLogin from '../../net/RespLogin';
import ReqLogin from '../../net/ReqLogin';
import { UIViewController } from '../../ui/UiViewController';
import { MainPaneController } from '../main/MainPaneController';
const { ccclass, property } = _decorator;

@ccclass('LoginPaneController')
export class LoginPaneController extends UIViewController {
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
    this.registerClickEvent(this.loginButton.node, this.onLoginClick, this);
  }

  onLoginClick() {
    const username = this.usernameInput.string;
    const password = this.passwordInput.string;

    // 这里添加登录验证逻辑
    if (username && password) {
      console.log('登录信息:', { username, password });

      GameContext.instance.WebSocketClient.sendMessage(
        ReqLogin.cmd,
        {
          id: username,
          pwd: password,
        },
        (msg: RespLogin) => {
          console.log('登录成功');
          MainPaneController.openUi();
        }
      );
    } else {
      console.log('请输入用户名和密码');
    }
  }
}
