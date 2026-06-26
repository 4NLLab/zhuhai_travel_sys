<script setup lang="ts">
import { onLaunch } from '@dcloudio/uni-app';
import { getRuntimeConfig } from '@/utils/config';
import { logger } from '@/adapters/logger';

onLaunch(() => {
  const config = getRuntimeConfig();
  logger.info('app_launch', {
    apiMode: config.apiMode,
    mockScenarioId: config.mockScenarioId,
    platform: config.platform
  });
});
</script>

<style>
page {
  min-height: 100%;
  background: var(--zt-color-bg);
  color: var(--zt-color-text);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}
</style>
