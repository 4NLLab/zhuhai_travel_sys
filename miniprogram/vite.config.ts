import { defineConfig, type PluginOption } from 'vite';
import uniPlugin from '@dcloudio/vite-plugin-uni';

const uni = (
  'default' in uniPlugin ? uniPlugin.default : uniPlugin
) as () => PluginOption | PluginOption[];

export default defineConfig({
  plugins: [uni()],
  resolve: {
    alias: {
      '@': '/src'
    }
  }
});
