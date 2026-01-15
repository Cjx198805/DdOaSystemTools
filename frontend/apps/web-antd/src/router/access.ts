import type {
  ComponentRecordType,
  GenerateMenuAndRoutesOptions,
} from '@vben/types';

import { generateAccessible } from '@vben/access';
import { preferences } from '@vben/preferences';

import { message } from 'ant-design-vue';

import { getAllMenusApi } from '#/api';
import { BasicLayout, IFrameView } from '#/layouts';
import { $t } from '#/locales';

const forbiddenComponent = () => import('#/views/_core/fallback/forbidden.vue');

async function generateAccess(options: GenerateMenuAndRoutesOptions) {
  console.log(
    '[Vben Access] generateAccess called with preferences.app.accessMode:',
    preferences.app.accessMode,
  );
  const pageMap: ComponentRecordType = import.meta.glob('../views/**/*.vue');
  console.log(
    '[Vben Access] pageMap keys (first 10):',
    Object.keys(pageMap).slice(0, 10),
  );

  const layoutMap: ComponentRecordType = {
    BasicLayout,
    IFrameView,
  };

  return await generateAccessible(preferences.app.accessMode, {
    ...options,
    fetchMenuListAsync: async () => {
      console.log('[Vben Access] fetchMenuListAsync triggered');
      message.loading({
        content: `${$t('common.loadingMenu')}...`,
        duration: 1.5,
      });
      const menuData = await getAllMenusApi();
      console.log(
        '[Vben Access] menuData from API:',
        JSON.stringify(menuData, null, 2),
      );

      // 递归标准化组件路径，确保匹配 pageMap
      const normalizeMenu = (menus: any[]) => {
        menus.forEach((menu) => {
          if (
            menu.component &&
            menu.component !== 'BasicLayout' &&
            menu.component !== 'IFrameView'
          ) {
            let comp = menu.component as string;
            if (comp.startsWith('/')) {
              comp = comp.slice(1);
            }
            if (!comp.endsWith('.vue')) {
              comp = `${comp}.vue`;
            }

            // 寻找匹配的 key
            const targetKey = Object.keys(pageMap).find(
              (key) => key.endsWith(comp) || key.includes(`/views/${comp}`),
            );

            if (targetKey) {
              console.log(
                `[Vben Access] Map component: ${menu.component} -> ${targetKey}`,
              );
              menu.component = targetKey;
            } else {
              console.warn(`[Vben Access] COMPONENT NOT FOUND: ${comp}`);
            }
          }

          // 强制生成安全的 name (Vue Router 不建议在 name 中包含斜杠)
          const cleanName =
            menu.path.replaceAll('/', '_').replace(/^_/, '') || 'Home';
          menu.name = cleanName;

          if (menu.children && menu.children.length > 0) {
            normalizeMenu(menu.children);
          }
        });
      };

      normalizeMenu(menuData);
      console.log(
        '[Vben Access] Final normalized menu:',
        JSON.stringify(menuData, null, 2),
      );
      return menuData;
    },
    // 可以指定没有权限跳转403页面
    forbiddenComponent,
    // 如果 route.meta.menuVisibleWithForbidden = true
    layoutMap,
    pageMap,
  });
}

export { generateAccess };
