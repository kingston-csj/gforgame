export class StringUtils {
  static isBlank(str: string): boolean {
    return str === null || str === undefined || str.trim() === '';
  }
}
