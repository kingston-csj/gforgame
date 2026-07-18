using System;
using System.Collections.Generic;
using Nova.Ui;
using UnityEngine;
using UnityEngine.Events;
using UnityEngine.UI;

namespace Nova.Mvc
{
    public abstract class BaseView : UiNode, IUiPanel
    {
        public PanelOpenIntend openIntend;

        // 核心成员
        protected BaseController _controller;

        // 事件容器：仅管理Button事件
        private readonly Dictionary<Button, List<UnityAction>> _buttonEvents =
            new Dictionary<Button, List<UnityAction>>();

        // IUIPanel 实现
        protected BaseView()
        {
        }

        public PanelState CurrentState { get; protected set; } = PanelState.None;

        #region IUIPanel 核心生命周期

        public void Init()
        {
            if (CurrentState != PanelState.None) return;

            CurrentState = PanelState.Created;
            try
            {
                InitUI(); // 子类初始化UI
            }
            catch (Exception e)
            {
                Debug.LogError($"Panel {GetType().Name} init failed: {e.Message}");
            }
        }

        public void Show(PanelOpenIntend openIntend, Action completeCallback)
        {
            this.openIntend = openIntend;
            if (CurrentState == PanelState.Showing) return;
            CurrentState = PanelState.Showing;
            try
            {
                gameObject.SetActive(true);
                OnShow(openIntend ?? new PanelOpenIntend());
                completeCallback?.Invoke();
            }
            catch (Exception e)
            {
                Debug.LogError($"Panel {GetType().Name} show failed: {e.Message}");
            }
        }

        public abstract void OnShow();

        public void Hide()
        {
            if (CurrentState == PanelState.Hiding || CurrentState == PanelState.Destroyed) return;

            CurrentState = PanelState.Hiding;
            try
            {
                gameObject.SetActive(false);
                OnHide();
            }
            catch (Exception e)
            {
                Debug.LogError($"Panel {GetType().Name} hide failed: {e.Message}");
            }
        }

        public void DestroyPanel()
        {
            if (CurrentState == PanelState.Destroyed) return;

            CurrentState = PanelState.Destroyed;
            try
            {
                UnregisterAllClickEvents(); // 销毁前清理所有事件
                OnDestroyPanel();
                Destroy(gameObject);
            }
            catch (Exception e)
            {
                Debug.LogError($"Panel {GetType().Name} destroy failed: {e.Message}");
            }
        }

        #endregion

        #region 核心：仅支持Button的事件注册/反注册（简洁版）

        /// <summary>
        /// 注册Button点击事件（仅支持Button，简洁无冗余）
        /// </summary>
        /// <param name="button">目标按钮</param>
        /// <param name="onClick">点击回调</param>
        protected void RegisterClickEvent(Button button, Action onClick)
        {
            if (button == null)
            {
                Debug.LogError($"[{GetType().Name}] RegisterClickEvent: Button is null!");
                return;
            }

            if (onClick == null)
            {
                Debug.LogError($"[{GetType().Name}] RegisterClickEvent: Callback is null!");
                return;
            }

            // 包装UnityAction
            UnityAction unityAction = () =>
            {
                try
                {
                    onClick.Invoke();
                }
                catch (Exception e)
                {
                    Debug.LogError($"[{GetType().Name}] Button {button.name} click error: {e.Message}");
                }
            };

            // 注册事件
            button.onClick.AddListener(unityAction);

            // 记录到容器
            if (!_buttonEvents.ContainsKey(button))
            {
                _buttonEvents[button] = new List<UnityAction>();
            }

            _buttonEvents[button].Add(unityAction);
        }

        /// <summary>
        /// 反注册指定Button的指定回调
        /// </summary>
        protected void UnregisterClickEvent(Button button, Action onClick)
        {
            if (button == null || onClick == null || !_buttonEvents.ContainsKey(button)) return;

            var targetActions = _buttonEvents[button];
            for (int i = targetActions.Count - 1; i >= 0; i--)
            {
                var unityAction = targetActions[i];
                button.onClick.RemoveListener(unityAction);
                targetActions.RemoveAt(i);
            }

            if (targetActions.Count == 0)
            {
                _buttonEvents.Remove(button);
            }
        }

        /// <summary>
        /// 反注册指定Button的所有事件
        /// </summary>
        protected void UnregisterClickEvent(Button button)
        {
            if (button == null || !_buttonEvents.ContainsKey(button)) return;

            var targetActions = _buttonEvents[button];
            foreach (var unityAction in targetActions)
            {
                button.onClick.RemoveListener(unityAction);
            }

            targetActions.Clear();
            _buttonEvents.Remove(button);
        }

        /// <summary>
        /// 反注册所有Button事件
        /// </summary>
        protected void UnregisterAllClickEvents()
        {
            foreach (var kvp in _buttonEvents)
            {
                var button = kvp.Key;
                var actions = kvp.Value;
                foreach (var action in actions)
                {
                    button.onClick.RemoveListener(action);
                }
            }

            _buttonEvents.Clear();
        }

        #endregion

        #region 子类扩展点（简洁版）

        // 必须重写：初始化UI组件
        protected abstract void InitUI();

        // 面板打开的显示逻辑
        protected virtual void OnShow(PanelOpenIntend openIntend)
        {
        }

        // 隐藏逻辑
        protected virtual void OnHide()
        {
        }


        // 可选重写：销毁清理逻辑
        protected virtual void OnDestroyPanel()
        {
        }

        #endregion

        #region Unity生命周期

        protected virtual void OnDestroy()
        {
            UnregisterAllClickEvents(); // 防止漏销毁
            _controller = null;
        }

        #endregion
    }
}
