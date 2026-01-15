<script lang="ts" setup>
import type { SystemFieldPermissionApi } from '#/api/system/field_permission';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { IconifyIcon } from '@vben/icons';

import { Button, message, Popconfirm, Space } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  deleteFieldPermission,
  getFieldPermissionList,
} from '#/api/system/field_permission';

import { useColumns, useGridFormSchema } from './data';
import PermissionForm from './modules/form.vue';

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  connectedComponent: PermissionForm,
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
          const res = await getFieldPermissionList({
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

function onEdit(row: SystemFieldPermissionApi.FieldPermission) {
  formDrawerApi.setData(row).open();
}

async function onDelete(row: SystemFieldPermissionApi.FieldPermission) {
  try {
    await deleteFieldPermission(row.id);
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
          新增字段权限
        </Button>
      </template>

      <template #operation="{ row }">
        <Space>
          <Button size="small" type="link" @click="onEdit(row)">
            <IconifyIcon icon="lucide:notebook-pen" class="size-4" />
            编辑
          </Button>
          <Popconfirm title="确定要删除该字段权限吗？" @confirm="onDelete(row)">
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
