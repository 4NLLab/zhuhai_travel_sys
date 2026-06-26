export const modalAdapter = {
  async confirm(message: string, title = '请确认'): Promise<boolean> {
    const result = await uni.showModal({
      title,
      content: message,
      confirmText: '确认',
      cancelText: '取消'
    });
    return result.confirm;
  },

  async alert(message: string, title = '提示'): Promise<void> {
    await uni.showModal({
      title,
      content: message,
      showCancel: false,
      confirmText: '知道了'
    });
  }
};
