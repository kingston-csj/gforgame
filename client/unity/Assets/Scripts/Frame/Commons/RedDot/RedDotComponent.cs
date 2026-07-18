namespace Nova.Commons.RedDot
{
    using UnityEngine;
    using UnityEngine.UI;

    /// <summary>
    /// 红点组件
    /// </summary>
    public class RedDotComponent : MonoBehaviour
    {
        /// <summary>
        /// 红点路径(展开为树型结构)
        /// </summary>
        public string path;

        public Text scoreTxt;

        /// <summary>
        /// 红点 ui节点
        /// </summary>
        public GameObject dot;

        /// <summary>
        /// 是否显示数字
        /// </summary>
        public bool showNum;

        /// <summary>
        /// 更新红点数字
        /// </summary>
        public void UpdateScore(int score) {
            this.dot.SetActive(score > 0);
            if (this.showNum) {
                this.scoreTxt.gameObject.SetActive(true);
                this.scoreTxt.text = score.ToString();
            }
        }

        protected void Awake()
        {
            if (!string.IsNullOrEmpty(this.path)) {
                this.dot = transform.Find(this.path).gameObject;
            }
            this.scoreTxt.gameObject.SetActive(false);
        }
    }
}
