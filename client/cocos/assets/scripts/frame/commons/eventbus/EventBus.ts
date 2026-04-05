/**
 * 事件总线工具类
 * 全局单例，用于组件间通信
 */
class EventBus {
    private static instance: EventBus;
    private eventListeners: Map<string, Array<(data?: any) => void>>;

    private constructor() {
        this.eventListeners = new Map();
    }

    /**
     * 获取事件总线单例
     */
    public static getInstance(): EventBus {
        if (!EventBus.instance) {
            EventBus.instance = new EventBus();
        }
        return EventBus.instance;
    }

    /**
     * 注册事件监听器
     * @param eventName 事件名称
     * @param callback 回调函数
     */
    public on(eventName: string, callback: (data?: any) => void): void {
        if (!this.eventListeners.has(eventName)) {
            this.eventListeners.set(eventName, []);
        }
        this.eventListeners.get(eventName)?.push(callback);
    }

    /**
     * 注册一次性事件监听器（触发后自动移除）
     * @param eventName 事件名称
     * @param callback 回调函数
     */
    public once(eventName: string, callback: (data?: any) => void): void {
        const onceCallback = (data?: any) => {
            callback(data);
            this.off(eventName, onceCallback);
        };
        this.on(eventName, onceCallback);
    }

    /**
     * 移除事件监听器
     * @param eventName 事件名称
     * @param callback 回调函数（可选，不提供则移除该事件的所有监听器）
     */
    public off(eventName: string, callback?: (data?: any) => void): void {
        if (!this.eventListeners.has(eventName)) return;

        if (callback) {
            const listeners = this.eventListeners.get(eventName);
            if (listeners) {
                const index = listeners.indexOf(callback);
                if (index > -1) {
                    listeners.splice(index, 1);
                }
            }
        } else {
            this.eventListeners.delete(eventName);
        }
    }

    /**
     * 触发事件
     * @param eventName 事件名称
     * @param data 事件数据（可选）
     */
    public emit(eventName: string, data?: any): void {
        if (!this.eventListeners.has(eventName)) return;

        const listeners = this.eventListeners.get(eventName);
        if (listeners) {
            listeners.forEach(callback => {
                try {
                    callback(data);
                } catch (error) {
                    console.error(`Error in event listener for ${eventName}:`, error);
                }
            });
        }
    }

    /**
     * 清空所有事件监听器
     */
    public clear(): void {
        this.eventListeners.clear();
    }

    /**
     * 获取指定事件的监听器数量
     * @param eventName 事件名称
     */
    public listenerCount(eventName: string): number {
        const listeners = this.eventListeners.get(eventName);
        return listeners ? listeners.length : 0;
    }
}

// 导出单例实例
export const eventBus = EventBus.getInstance();
export default eventBus;
