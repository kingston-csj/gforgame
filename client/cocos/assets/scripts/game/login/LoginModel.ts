import GameContext from '../../GameContext';
import ReqLogin from '../../net/protocol/ReqLogin';
import RespLogin from '../../net/protocol/RespLogin';

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
    return new Promise<RespLogin>((resolve, reject) => {
      GameContext.instance.WebSocketClient.sendMessage(
        ReqLogin.cmd,
        {
          id: this.userId,
          pwd: this.userPwd,
        },
        (msg: RespLogin) => {
          resolve(msg);
        }
      );
    });
  }
}
