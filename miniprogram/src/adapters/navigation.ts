export const navigationAdapter = {
  switchTab(url: string): void {
    uni.switchTab({ url });
  },

  navigateTo(url: string): void {
    uni.navigateTo({ url });
  },

  back(delta = 1): void {
    uni.navigateBack({ delta });
  }
};
