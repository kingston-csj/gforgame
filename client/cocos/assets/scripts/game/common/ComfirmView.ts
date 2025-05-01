import { _decorator, Button, Label } from 'cc';
import { BaseUiView } from '../../ui/BaseUiView';

const { ccclass, property } = _decorator;

@ccclass('ComfirmView')
export class ComfirmView extends BaseUiView {
  @property(Label)
  titleLabel: Label = null!;

  @property(Label)
  contentLabel: Label = null!;

  @property(Button)
  confirmButton: Button = null!;

  @property(Button)
  cancelButton: Button = null!;

  public setTitle(title: string) {
    this.titleLabel.string = title;
  }

  public setContent(content: string) {
    this.contentLabel.string = content;
  }
}
