using Frame.Commons.Utils;
using Nova.Mvc;

namespace Game.Login
{
    public class LoginController : BaseController
    {
        // 绑定 View 事件
        protected override void BindViewEvents()
        {
            LoginView _loginView = _view as LoginView;
            LoginModel _loginModel = _model as LoginModel;
            // 账号输入变更 → 更新 Model
            _loginView.OnUsernameChanged += (username) => { _loginModel.Username = username; };
            // 登录按钮点击 → 调用 Model 处理登录
            _loginView.OnLoginClicked += OnLoginClicked;
        }

        // 处理登录逻辑
        private void OnLoginClicked()
        {
            LoginModel _loginModel = _model as LoginModel;
            LoginView _loginView = _view as LoginView;
            if (_loginModel.Username.IsEmpty())
            {
                _loginView.UpdateTipText("请输入账号密码登录");
                return;
            }
            _loginModel.DoLogin();
        }
    }
}