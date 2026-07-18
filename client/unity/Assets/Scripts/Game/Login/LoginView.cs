using System;
using Nova.Mvc;
using UnityEngine;
using UnityEngine.UI;
using Button = UnityEngine.UI.Button;

namespace Game.Login
{
    // 登录视图：负责UI显示和交互
    public class LoginView : BaseView
    {
        [Header("登录UI组件")] [SerializeField] private InputField _usernameInput;
        [SerializeField] private InputField _passwordInput;
        [SerializeField] private Button _loginBtn;
        [SerializeField] private Text _tipText;

        public event Action<string> OnUsernameChanged;

        // 密码输入变更事件
        public event Action<string> OnPasswordChanged;

        // 登录按钮点击事件
        public event Action OnLoginClicked;


        // 初始化UI组件
        protected override void InitUI()
        {
            // 绑定输入框变更事件
            _usernameInput.onValueChanged.AddListener(OnUsernameInputChanged);
            _passwordInput.onValueChanged.AddListener(OnPasswordInputChanged);

            // 绑定按钮点击事件
            RegisterClickEvent(_loginBtn, OnLoginBtnClicked);

            // 初始化提示文本
            UpdateTipText("请输入账号密码登录");

            LoginController controller = new LoginController();
            controller.Init(this, LoginModel.Instance);
        }

        // 更新用户名显示
        public void UpdateUsername(string username)
        {
            if (_usernameInput != null && _usernameInput.text != username)
            {
                _usernameInput.text = username;
            }
        }

        public override void OnShow()
        {
            
        }

        // 更新密码显示
        public void UpdatePassword(string password)
        {
            if (_passwordInput != null && _passwordInput.text != password)
            {
                _passwordInput.text = password;
            }
        }

        // 更新提示文本
        public void UpdateTipText(string text)
        {
            if (_tipText != null)
            {
                _tipText.text = text;
            }
        }

        // 清空输入框
        public void ClearInputs()
        {
            UpdateUsername(string.Empty);
            UpdatePassword(string.Empty);
            UpdateTipText("请输入账号密码登录");
        }

        #region 内部事件处理

        private void OnUsernameInputChanged(string value)
        {
            OnUsernameChanged?.Invoke(value);
        }

        private void OnPasswordInputChanged(string value)
        {
            OnPasswordChanged?.Invoke(value);
        }

        private void OnLoginBtnClicked()
        {
            OnLoginClicked?.Invoke();
        }

        #endregion

        // 销毁时清理事件
        protected override void OnDestroy()
        {
            base.OnDestroy();

            if (_usernameInput != null) _usernameInput.onValueChanged.RemoveAllListeners();
            if (_passwordInput != null) _passwordInput.onValueChanged.RemoveAllListeners();

            if (_loginBtn != null) UnregisterClickEvent(_loginBtn, OnLoginBtnClicked);

            OnUsernameChanged = null;
            OnPasswordChanged = null;
            OnLoginClicked = null;
        }

        public void Dispose()
        {
            // TODO release managed resources here
        }
    }
}