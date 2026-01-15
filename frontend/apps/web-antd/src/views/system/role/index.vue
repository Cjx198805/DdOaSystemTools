<script lang="ts" setup>
import type { SystemRoleApi } from '#/api/system/role';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { IconifyIcon } from '@vben/icons';

import { Button, message, Popconfirm, Space } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { deleteSystemRole, getSystemRoleList } from '#/api/system/role';

import { useColumns, useGridFormSchema } from './data';
import RoleForm from './modules/form.vue';

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  connectedComponent: RoleForm,
});

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: {
    schema: useGridFormSchema(),
    submitOnChange: true,
  },
  gridOptions: {
    columns: useColumns(),
    height: 'auto',
    proxyConfig: {
      ajax: {
        query: async ({ page }, formValues) => {
          const res = await getSystemRoleList({
            page: page.currentPage,
            page_size: page.pageSize,
            ...formValues,
          });
          return res;
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

function onCreate() {
  formDrawerApi.open();
}

function onEdit(row: SystemRoleApi.Role) {
  formDrawerApi.setData(row).open();
}

async function onDelete(row: SystemRoleApi.Role) {
  try {
    await deleteSystemRole(row.id);
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
        <Button type="primary" @click="onCreate">
          <IconifyIcon icon="lucide:plus" class="mr-1 size-4" />
          新增角色
        </Button>
      </template>

      <template #operation="{ row }">
        <Space>
          <Button size="small" type="link" @click="onEdit(row)">
            <IconifyIcon icon="lucide:notebook-pen" class="size-4" />
            编辑
          </Button>
          <Popconfirm title="确定要删除该角色吗？" @confirm="onDelete(row)">
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
