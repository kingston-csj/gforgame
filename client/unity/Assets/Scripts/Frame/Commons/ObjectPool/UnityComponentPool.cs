namespace Frame.Commons.ObjectPool
{
    using UnityEngine;

    /// <summary>
    /// Unity预制体专用对象池（基于Component，适配业务组件复用）
    /// </summary>
    /// <typeparam name="T">限制为Component（如MonoBehaviour子类，绑定业务逻辑）</typeparam>
    public class UnityComponentPool<T> : GenericObjectPool<T> where T : Component
    {
        #region 核心参数

        private readonly GameObject _prefab; // 预制体，池化对象的模板
        private readonly Transform root; // 池对象父节点（层级管理）

        #endregion

        #region 构造函数

        /// <summary>
        /// 初始化Component类型的预制体池
        /// </summary>
        /// <param name="prefab">预制体（需挂载T类型的Component，不能为空）</param>
        /// <param name="root">池父节点（可选，自动创建全局节点）</param>
        public UnityComponentPool(GameObject prefab, Transform root = null)
        {
            _prefab = prefab ?? throw new System.ArgumentNullException(nameof(prefab), "预制体不能为空，且必须挂载T类型的Component");
            this.root = root;
        }

        #endregion

        #region 重写父类核心方法（适配Component）

        /// <summary>
        /// 创建新对象：实例化预制体并返回绑定的Component
        /// </summary>
        protected override T CreateNewObject()
        {
            // 实例化预制体（基于Component的预制体，Instantiate后直接返回对应Component）
            GameObject go = GameObject.Instantiate(_prefab, root);
            go.transform.localPosition = Vector3.zero;
            return go.GetComponent<T>();
        }

        /// <summary>
        /// 获取对象时：激活、重置状态、解绑父节点
        /// </summary>
        protected override void OnItemAcquired(T item)
        {
            base.OnItemAcquired(item);
            // 核心操作：激活游戏对象
            item.gameObject.SetActive(true);
        }

        protected override void OnItemReturned(T item)
        {
            base.OnItemReturned(item);
            if (item == null) return;
            // 核心操作：禁用游戏对象 + 归位池父节点
            item.gameObject.SetActive(false);
            item.transform.SetParent(root);
        }

        #endregion

        #region 扩展方法（Unity业务场景专用）

        /// <summary>
        /// 安全归还对象
        /// </summary>
        public new void Return(T obj)
        {
            if (obj == null) return;
            OnItemReturned(obj); // 先执行回收逻辑
            base.Return(obj); // 再加入池队列
        }

        #endregion
    }
}