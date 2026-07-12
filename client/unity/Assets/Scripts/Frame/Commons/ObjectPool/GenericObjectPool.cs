namespace Frame.Commons.ObjectPool
{
    using System;
    using System.Collections.Generic;

    /// <summary>
    /// 通用对象池（无Unity依赖、无容量限制、无需回收、高度通用）
    /// 支持任意引用类型的对象复用，避免频繁new/GC
    /// </summary>
    /// <typeparam name="T">要池化的对象类型（必须是引用类型）</typeparam>
    public class GenericObjectPool<T> where T : class
    {
        private readonly Queue<T> _pool; // 空闲对象队列

        #region 构造函数

        /// <summary>
        /// 初始化通用对象池
        /// </summary>
        public GenericObjectPool()
        {
            _pool = new Queue<T>();
        }

        #endregion

        #region 核心方法

        /// <summary>
        /// 从池中获取对象（无空闲则创建新对象）
        /// </summary>
        /// <returns>复用/新建的对象</returns>
        public T Borrow()
        {
            T obj;
            // 有空闲对象则复用，无则创建新对象
            if (_pool.Count > 0)
            {
                obj = _pool.Dequeue();
            }
            else
            {
                obj = CreateNewObject();
            }

            OnItemAcquired(obj);

            return obj;
        }

        public void Return(T obj)
        {
            _pool.Enqueue(obj);
        }

        /// <summary>
        /// 清空对象池（释放所有空闲对象，触发GC）
        /// </summary>
        public void Clear()
        {
            _pool.Clear();
        }

        #endregion

        #region 私有方法

        /// <summary>
        /// 创建新对象
        /// </summary>
        protected virtual T CreateNewObject()
        {
            return default;
        }

        /// <summary>
        /// 对象获取触发点
        /// </summary>
        /// <param name="item"></param>
        protected virtual void OnItemAcquired(T item)
        {
        }

        /// <summary>
        /// 对象归还触发点
        /// </summary>
        /// <param name="item"></param>
        protected virtual void OnItemReturned(T item)
        {
        }

        #endregion
    }
}