using Game.Core;
using Game.Net.Message;
using Nova.Mvc;
using UnityEngine;

namespace Game.Login
{
// 登录模型：处理登录数据和业务逻辑
    public class LoginModel : BaseModel
    {
        // 定义字段名常量（避免魔法字符串）
        public const string FIELD_USERNAME = "Username";
        public const string FIELD_LOGIN_RESULT = "LoginResult";

        private string _username;
        private string _password;
        private bool _loginSuccess;

        private static LoginModel _instance;

        public static LoginModel Instance
        {
            get
            {
                if (_instance == null)
                    _instance = new LoginModel();
                return _instance;
            }
        }

        // 用户名（变更时通知）
        public string Username
        {
            get => _username;
            set
            {
                if (_username != value)
                {
                    _username = value;
                    Notify(FIELD_USERNAME, value);
                }
            }
        }

        // 密码（变更时通知）
        public string Password { get; set; }

        // 登录结果（变更时通知）
        public bool LoginSuccess
        {
            get => _loginSuccess;
            set
            {
                if (_loginSuccess != value)
                {
                    _loginSuccess = value;
                    Notify(FIELD_LOGIN_RESULT, value);
                }
            }
        }

        // 登录验证
        public void DoLogin()
        {
            // 异步登录
            Debug.Log("开始验证登录信息...");
            WebManager.Instance.ConnectToServer(AppContext.gameConfig.serverUrl, () =>
            {
                Debug.Log("登录验证成功！");
                // 发送登录请求
                ReqPlayerLogin reqLogin = new ReqPlayerLogin { playerId = _username };
                AppContext.socketClient.Send(reqLogin, (ResPlayerLogin res) => Debug.Log($"登录成功，玩家名称：{res.name}"));
            });
        }

        // 清空登录信息
        public void ClearLoginInfo()
        {
            Username = string.Empty;
            Password = string.Empty;
            LoginSuccess = false;
        }

        public override void ClearCallbacks()
        {
            base.ClearCallbacks();
            ClearLoginInfo();
        }
    }
}