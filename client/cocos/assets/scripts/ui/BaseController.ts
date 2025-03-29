import { Component } from 'cc';

export class BaseController extends Component {
  protected view: any;

  public initView(view: any) {
    this.view = view;
    // 进行一些通用的视图初始化操作，比如设置视图的父节点等
    this.view.node.parent = this.node;
    // 调用绑定事件的方法
    this.bindViewEvents();
  }

  protected bindViewEvents() {
    // 子类可以重写此方法来绑定具体的事件
  }
}
