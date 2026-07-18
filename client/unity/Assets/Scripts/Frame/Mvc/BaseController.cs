namespace Nova.Mvc
{
    public class BaseController
    {
        protected BaseView _view;
        protected BaseModel _model;

        // 初始化控制器（关联View和Model）
        public virtual void Init(BaseView view, BaseModel model)
        {
            _view = view;
            _model = model;
            _view.Init();
            InitView(view);
            BindViewEvents();
            BindModelEvents();
        }

        public void InitView(BaseView view)
        {
            _view = view;
        }

        // 绑定View事件（子类重写）
        protected virtual void BindViewEvents()
        {
        }

        // 绑定Model事件（子类重写）
        protected virtual void BindModelEvents()
        {
        }

        // 销毁时清理引用
        public virtual void Dispose()
        {
            _view = null;
            _model = null;
        }
    }
}