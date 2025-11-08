namespace Nova.Mvc
{
    
    public class BaseController
    {
        protected BaseView _view;

        public void InitView(BaseView view)
        {
            _view = view;
        }

        protected void BindViewEvents()
        {
            // 子类可以重写此方法来绑定具体的事件
        }
    }
}