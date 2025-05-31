import { _decorator } from 'cc';
import PlayerData from '../../data/user/PlayerData';
import { BaseController } from '../../frame/mvc/BaseController';
import { TipsPaneController } from '../common/TipsPaneController';
import GameConstants from '../constants/GameConstants';
import { MainPaneController } from '../main/MainPaneController';
import { LoginModel } from './LoginModel';
import { LoginPaneView } from './LoginPaneView';
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
        if (response.code === 0) {
          console.log('登录成功');
          PlayerData.instance.name = response.name;
          PlayerData.instance.fighting = response.fighting;
          PlayerData.instance.camp = response.camp;
          MainPaneController.openUi();
        } else {
          TipsPaneController.showI18nContent(response.code);
        }
      } catch (error) {
        console.error('登录失败:', error);
        TipsPaneController.showI18nContent(GameConstants.I18N.ILLEGAL_PARAMS);
      }
    } else {
      console.log('请输入用户名和密码');
      TipsPaneController.showI18nContent(GameConstants.I18N.CONTENT_NOT_ENOUGH);
    }
  }
}
