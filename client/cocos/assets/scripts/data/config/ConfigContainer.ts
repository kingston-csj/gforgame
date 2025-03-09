export default class ConfigContainer {
  private static _instance: ConfigContainer;
  private _configMap: Map<string, any> = new Map();

  private constructor() {}
}
