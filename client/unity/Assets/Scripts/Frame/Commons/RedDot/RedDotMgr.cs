﻿namespace Nova.Commons.RedDot
{
    using System.Collections.Generic;
    using UnityEngine;

    /// <summary>
    /// 红点管理器
    /// </summary>
    public class RedDotMgr
    {
        /// <summary>
        /// 根节点
        /// </summary>
        private RedDotNode root = new RedDotNode();

        private static RedDotMgr instance = new RedDotMgr();

        /// <summary>
        /// 绑定红点
        /// <param name="path">红点路径</param>
        /// <param name="ui">红点UI</param>
        /// </summary>
        public void Binding(string path, RedDotComponent ui)
        {
            var node = this.GetOrCreateNode(path);
            node.ui = ui;
            node.score = node.score;
        }

        private RedDotNode GetOrCreateNode(string path)
        {
            var node = this.root;
            foreach (var key in path.Split('/'))
            {
                node = node.AddChild(key);
            }
            return node;
        }

        /// <summary>
        /// 更新红点分数
        /// <param name="path">红点路径</param>
        /// <param name="score">红点分数</param>
        /// </summary>
        public void UpdateScore(string path, int score)
        {
            var node = this.GetOrCreateNode(path);
            node.score = score;
             // 向上回溯，更新父节点分数
             var parent = node.parent;
             while (parent != null)
             {
                int parentScore = 0;
                foreach (var child in parent.children.Values)
                {
                    parentScore += child.score;
                }
                 parent.score = parentScore;
                 parent = parent.parent;
             }
             // 从当前节点开始，向上更新UI
             var curr = node;
             while (curr != null)
             {
                curr.ui?.UpdateScore(curr.score);
                curr = curr.parent;
             }
        }
    }
}
