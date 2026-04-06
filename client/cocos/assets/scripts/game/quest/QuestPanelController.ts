import { _decorator } from 'cc';
import { BaseController } from "../../frame/mvc/BaseController";
import QuestView from './QuestView';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
import { LayerIdx } from '../../ui/LayerIds';
const { ccclass, property } = _decorator;


@ccclass('QuestPanelController')
export default class QuestPanelController extends BaseController {
  
    private static instance: QuestPanelController;
    private static creatingPromise: Promise<QuestPanelController> | null = null;

    @property(QuestView)
    public view: QuestView;

    public static openUi() {
    this.getInstance().then((controller) => {
        controller.view.display();
    });
    }

    protected start(): void {
    this.initView(this.view);
    }

    private static getInstance(): Promise<QuestPanelController> {
    if (this.instance) {
        return Promise.resolve(this.instance);
    }

    if (this.creatingPromise) {
        return this.creatingPromise;
    }
    this.creatingPromise = new Promise((resolve) => {
        UiViewFactory.createUi(R.Prefabs.QuestPanel, LayerIdx.layer4, (ui: QuestPanelController) => {
        this.instance = ui;
        this.creatingPromise = null;
        resolve(ui);
        });
    });
    return this.creatingPromise;
    }
}