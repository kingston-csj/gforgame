import { _decorator, EventTouch, Label, Node, Vec2, Vec3 } from 'cc';
import { ConfigContext } from '../../data/config/container/ConfigContext';
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import { HeroVo } from '../../net/protocol/items/HeroVo';
import AssetResourceFactory from '../../ui/AssetResourceFactory';
import R from '../../ui/R';
import { UiUtil } from '../../utils/UiUtil';
import { TipsPaneController } from '../common/TipsPaneController';
import { HeroBoxModel } from '../hero/HeroBoxModel';
import { BuZhenPaneController } from './BuZhenPaneController';
const { ccclass, property } = _decorator;

@ccclass('BuZhenHeroUpItem')
export class BuZhenHeroUpItem extends BaseUiView {
  @property(Label)
  private heroName: Label;

  @property(Node)
  private icon: Node;

  @property(Node)
  private camp: Node;

  @property(Label)
  private level: Label;

  private touchStartPos: Vec2 = new Vec2();
  private nodeStartPos: Vec3 = new Vec3();
  private hero: HeroVo;

  private allPositions: Node[] = [];

  public fillData(hero: HeroVo) {
    this.hero = hero;

    let heroData = ConfigContext.configHeroContainer.getRecord(hero.id);
    this.heroName.string = heroData.name;
    let heroSpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Hero);
    UiUtil.fillSpriteContent(this.icon, heroSpriteAtlas.getSpriteFrame(heroData.icon));
    let campSpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Camp);
    UiUtil.fillSpriteContent(this.camp, campSpriteAtlas.getSpriteFrame('camp_' + heroData.camp));
    this.level.string = hero.level.toString();
  }

  start() {
    this.node.on(Node.EventType.TOUCH_START, this.onTouchStart, this);
    this.node.on(Node.EventType.TOUCH_MOVE, this.onTouchMove, this);
    this.node.on(Node.EventType.TOUCH_END, this.onTouchEnd, this);

    // 获取所有可放置的位置
    this.allPositions = [
      this.node.parent.parent.parent.getChildByName('pos1'),
      this.node.parent.parent.parent.getChildByName('pos2'),
      this.node.parent.parent.parent.getChildByName('pos3'),
      this.node.parent.parent.parent.getChildByName('pos4'),
      this.node.parent.parent.parent.getChildByName('pos5'),
    ];
  }

  onTouchStart(event: EventTouch) {
    const location = event.getUILocation();
    this.touchStartPos.set(location.x, location.y);
    this.nodeStartPos.set(this.node.position.x, this.node.position.y, 0);
  }

  onTouchMove(event: EventTouch) {
    const location = event.getUILocation();
    const deltaX = location.x - this.touchStartPos.x;
    const deltaY = location.y - this.touchStartPos.y;

    this.node.setPosition(this.nodeStartPos.x + deltaX, this.nodeStartPos.y + deltaY);
  }

  onTouchEnd() {
    let nearestPos = this.findNearestPosition();
    // 如果找到最近的位置，且距离小于100，则移动到该位置
    if (nearestPos) {
      const posIndex = this.allPositions.indexOf(nearestPos) + 1;
      if (posIndex != this.hero.position) {
        HeroBoxModel.getInstance()
          .requestChangePostion(this.hero.id, posIndex)
          .then((msg) => {
            if (msg.code == 0) {
              if (msg.heroA > 0) {
                HeroBoxModel.getInstance().getHero(msg.heroA).position = msg.posA;
              }
              if (msg.heroB > 0) {
                HeroBoxModel.getInstance().getHero(msg.heroB).position = msg.posB;
              }
              BuZhenPaneController.refreshLineupHeros();
            } else {
              TipsPaneController.showI18nContent(msg.code);
            }
          });
      }
    } else {
      // 如果距离太远，回到原始位置
      this.node.setPosition(this.nodeStartPos);
    }
  }

  private findNearestPosition(): Node {
    let minDistance = Number.MAX_VALUE;
    let nearestPos: Node = null;

    // 获取当前节点的世界坐标
    const currentWorldPos = this.node.getWorldPosition();

    // 遍历所有位置，找到最近的
    this.allPositions.forEach((pos) => {
      if (pos) {
        const posWorldPos = pos.getWorldPosition();
        const distance = Vec3.distance(currentWorldPos, posWorldPos);

        if (distance < minDistance) {
          minDistance = distance;
          nearestPos = pos;
        }
      }
    });

    // 如果找到最近的位置，且距离小于100，则符合目标
    if (nearestPos && minDistance < 100) {
      return nearestPos;
    }
    return null;
  }
}
