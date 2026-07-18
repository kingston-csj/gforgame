using System;
using Nova.Mvc;

namespace Nova.Ui
{
    public interface IUiPanel
    {
        PanelState CurrentState { get; } // 当前面板状态

        /// <summary>
        /// 初始化面板
        /// </summary>
        void Init();

        /// <summary>
        /// 显示面板
        /// </summary>
        /// <param name="openIntend"></param>
        void Show(PanelOpenIntend openIntend, Action completeCallback);

        void OnShow();

        /// <summary>
        /// 隐藏面板
        /// </summary>
        void Hide();

        void DestroyPanel();
    }
}