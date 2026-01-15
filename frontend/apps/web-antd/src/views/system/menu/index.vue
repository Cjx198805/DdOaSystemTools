<script lang="ts" setup>
import type { SystemMenuApi } from '#/api/system/menu';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { IconifyIcon } from '@vben/icons';

import { Button, message, Popconfirm, Space } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { deleteSystemMenu, getSystemMenuTree } from '#/api/system/menu';

import { useColumns } from './data';
import MenuForm from './modules/form.vue';

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  connectedComponent: MenuForm,
});

const [Grid, gridApi] = useVbenVxeGrid({
  gridOptions: {
    columns: useColumns(),
    height: 'auto',
    treeConfig: {
      transform: false, // 后端已经返回了 tree 结构
      rowField: 'id',
      parentField: 'parent_id',
    },
    proxyConfig: {
      ajax: {
        query: async () => {
          const res = await getSystemMenuTree();
          return { items: res }; // Vben expects { items: [] } for tree usually
        },
      },
    },
    rowConfig: {
      keyField: 'id',
    },
    toolbarConfig: {
      custom: true,
      refresh: true,
      zoom: true,
    },
  },
});

function onRefresh() {
  gridApi.query();
}

function onCreate(row?: SystemMenuApi.Menu) {
  formDrawerApi.setData({ parent_id: row?.id || 0 }).open();
}

function onEdit(row: SystemMenuApi.Menu) {
  formDrawerApi.setData(row).open();
}

async function onDelete(row: SystemMenuApi.Menu) {
  try {
    await deleteSystemMenu(row.id);
    message.success('删除成功');
    onRefresh();
  } catch {}
}
</script>

<template>
  <Page auto-content-height>
    <FormDrawer @success="onRefresh" />
    <Grid>
      <template #toolbar-tools>
        <Button type="primary" @click="onCreate()">
          <IconifyIcon icon="lucide:plus" class="mr-1 size-4" />
          新增顶级菜单
        </Button>
      </template>

      <template #operation="{ row }">
        <Space>
          <Button size="small" type="link" @click="onCreate(row)">
            新增子项
          </Button>
          <Button size="small" type="link" @click="onEdit(row)">
            <IconifyIcon icon="lucide:notebook-pen" class="size-4" />
            编辑
          </Button>
          <Popconfirm
            title="确定要删除该菜单及其子项吗？"
            @confirm="onDelete(row)"
          >
            <Button danger size="small" type="link">
              <IconifyIcon icon="lucide:trash-2" class="size-4" />
              删除
            </Button>
          </Popconfirm>
        </Space>
      </template>
    </Grid>
  </Page>
</template>
