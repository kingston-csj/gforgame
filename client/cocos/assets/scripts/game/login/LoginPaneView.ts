import { _decorator, Component, Node, EditBox, Button, director } from 'cc';
import { BaseUiView } from '../../ui/BaseUiView';

const { ccclass, property } = _decorator;

@ccclass('LoginPaneView')
export class LoginPaneView extends BaseUiView {
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
    this.registerClickEvent(
      this.loginButton.node,
      () => {
        this.node.emit('accountLogin');
      },
      this
    );
  }

  public getUserId(): string {
    return this.usernameInput.string;
  }

  public getUserPwd(): string {
    return this.passwordInput.string;
  }
}
