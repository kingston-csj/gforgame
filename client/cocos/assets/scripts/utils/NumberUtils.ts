export class NumberUtils {
  // 超过1万，显示11.11万， 不足1万显示全部
  public static formatNumber(num: number): string {
    if (num >= 10000) {
      return (num / 10000).toFixed(2) + '万';
    } else {
      return num.toString();
    }
  }
}
