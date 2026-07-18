﻿namespace Nova.Commons.RedDot
{
    using System.Collections.Generic;
    using UnityEngine;

    /// <summary>
    /// 红点节点
    /// </summary>
    public class RedDotNode  
    {
       
       public int score {
           get; set;
       }


        public RedDotComponent ui  {
           get;  set;
        }


        /// <summary>
        /// 父节点
        /// </summary>
       public RedDotNode parent{
        get  ; set;
       }

          /// <summary>
          /// 构造函数
          /// </summary>
          /// <param name="score">初始分数</param>
        public RedDotNode(int score = 0) {
            this.score = score;
            this.children = new Dictionary<string, RedDotNode>();
        }
       /// <summary>
       /// 子节点
       /// </summary>
       public Dictionary<string, RedDotNode> children  {
        get; private set;
       }
        
        public Dictionary<string, RedDotNode> GetChildren()
        {
            return this.children;
        }
        public RedDotNode AddChild(string key)  {
           if (this.children.TryGetValue(key, out var child)) {
               return child;
           }
           child = new RedDotNode();
           child.parent = this;
           this.children[key] = child;
           return child;
        }
    }
}
