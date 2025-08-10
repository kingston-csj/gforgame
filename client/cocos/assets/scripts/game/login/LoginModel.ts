import GameContext from '../../GameContext';
import ReqLogin from '../../net/protocol/ReqLogin';
import RespLogin from '../../net/protocol/RespLogin';
import { StringUtils } from '../../utils/StringUtils';

export class LoginModel {
  private static _instance: LoginModel;

  private userId: string = '';

  private userPwd: string = '';

  public static get instance(): LoginModel {
    if (!LoginModel._instance) {
      LoginModel._instance = new LoginModel();
    }
    return LoginModel._instance;
  }

  public setUserId(userId: string) {
    this.userId = userId;
  }

  public setUserPwd(userPwd: string) {
    this.userPwd = userPwd;
  }

  public getUserId(): string {
    return this.userId;
  }

  public getUserPwd(): string {
    return this.userPwd;
  }

  public login(): Promise<RespLogin> {
    if (StringUtils.isBlank(this.userId) || StringUtils.isBlank(this.userPwd)) {
      return Promise.reject(new Error('用户名或密码不能为空'));
    }
    return new Promise<RespLogin>((resolve, reject) => {
      GameContext.instance.WebSocketClient.sendMessage(
        ReqLogin.cmd,
        {
          PlayerId: this.userId,
          pwd: this.userPwd,
        },
        (msg: RespLogin) => {
          resolve(msg);
        }
      );
    });
  }
}
