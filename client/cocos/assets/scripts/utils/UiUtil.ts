import { Node, Sprite, SpriteFrame, UITransform } from 'cc';

export class UiUtil {
  public static fillSpriteContent(icon: Node, spriteFrame: SpriteFrame) {
    const iconTransform = icon.getComponent(UITransform);
    if (!iconTransform) {
      console.warn('Icon node has no UITransform component');
      return;
    }
    // 保存节点当前的尺寸，用于调整图像
    const originalIconWidth = iconTransform.width;
    const originalIconHeight = iconTransform.height;

    icon.getComponent(Sprite).spriteFrame = spriteFrame;
    // 获取当前SpriteFrame
    const sprite = icon.getComponent(Sprite);
    if (!sprite || !sprite.spriteFrame) {
      console.warn('Icon has no valid sprite frame');
      return;
    }
    // 设置UITransform的contentSize为原始图片尺寸
    iconTransform.setContentSize(originalIconWidth, originalIconHeight);
  }
}
