export class BaseModel {
  private changeCallbacks: Map<string, ((value: any) => void)[]> = new Map();

  // Register change listener for any field
  onChange(field: string, callback: (value: any) => void) {
    if (!this.changeCallbacks.has(field)) {
      this.changeCallbacks.set(field, []);
    }
    this.changeCallbacks.get(field)!.push(callback);
    return () => {
      const callbacks = this.changeCallbacks.get(field);
      if (callbacks) {
        const index = callbacks.indexOf(callback);
        if (index > -1) {
          callbacks.splice(index, 1);
        }
      }
    };
  }

  // Notify change for any field
  protected notifyChange(field: string, value: any) {
    const callbacks = this.changeCallbacks.get(field);
    if (callbacks) {
      callbacks.forEach((callback) => callback(value));
    }
  }
}
