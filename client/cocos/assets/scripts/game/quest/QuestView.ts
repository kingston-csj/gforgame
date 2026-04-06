import { _decorator, instantiate, Node, Prefab } from 'cc';
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import { QuestBoxModel } from './QuestBoxModel';
import GameConstants from '../constants/GameConstants';
import { QuestItemView } from './QuestItemView';
const { ccclass, property } = _decorator;


@ccclass('QuestView')
export default class QuestView extends BaseUiView {
    @property(Node)
    questContainer: Node;

    @property(Prefab)
    questItemPrefab: Prefab;

    @property(Node)
    closeBtn: Node;

    protected start(): void {
         this.registerClickEvent(this.closeBtn, this.hide, this);
     }

    protected onDisplay(): void {
        this.questContainer.children.forEach((child) => {
            child.destroy();
        });
        const quests = QuestBoxModel.instance.getQuestsByCategory(GameConstants.Quest.Category.DAILY);
        for (const quest of quests) {
            const questItem = instantiate(this.questItemPrefab);
            questItem.setParent(this.questContainer);
            questItem.getComponent(QuestItemView).fillData(quest);
        }
    }
}