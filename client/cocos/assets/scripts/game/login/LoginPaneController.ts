import { _decorator, Button, Component } from 'cc';
import { MainPaneController } from '../main/MainPaneController';
import { LoginPaneView } from './LoginPaneView';
import { LoginModel } from './LoginModel';
import { BaseController } from '../../ui/BaseController';
const { ccclass, property } = _decorator;

@ccclass('LoginPaneController')
export class LoginPaneController extends BaseController {
  @property(LoginPaneView)
  loginView: LoginPaneView | null = null;

  loginModel: LoginModel = LoginModel.instance;

  start() {
    this.initView(this.loginView);
  }

  protected bindViewEvents() {
    this.loginView.node.on('accountLogin', this.onLoginClick, this);
  }

  async onLoginClick() {
    const username = this.loginView.getUserId();
    const password = this.loginView.getUserPwd();

    // 这里添加登录验证逻辑
    if (username && password) {
      this.loginModel.setUserId(username);
      this.loginModel.setUserPwd(password);

      try {
        const response = await this.loginModel.login();
        console.log('登录成功');
        MainPaneController.openUi();
      } catch (error) {
        console.error('登录失败:', error);
      }
    } else {
      console.log('请输入用户名和密码');
    }
  }
}
