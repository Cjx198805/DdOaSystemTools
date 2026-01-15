<script lang="ts" setup>
import type { BusinessDownloadApi } from '#/api/business/download';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { IconifyIcon } from '@vben/icons';

import { Button, message, Popconfirm, Progress, Space } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  deleteDownloadTask,
  getDownloadTaskList,
} from '#/api/business/download';

import { useColumns, useGridFormSchema } from './data';
import DownloadForm from './modules/form.vue';

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  connectedComponent: DownloadForm,
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
          const res = await getDownloadTaskList({
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

async function onDelete(row: BusinessDownloadApi.Task) {
  try {
    await deleteDownloadTask(row.id);
    message.success('删除成功');
    onRefresh();
  } catch {}
}

function onDownloadResult(row: BusinessDownloadApi.Task) {
  if (row.result_url) {
    window.open(row.result_url, '_blank');
  } else {
    message.warning('结果文件尚未生成');
  }
}
</script>

<template>
  <Page auto-content-height>
    <FormDrawer @success="onRefresh" />
    <Grid>
      <template #toolbar-tools>
        <Button type="primary" @click="onCreate">
          <IconifyIcon icon="lucide:plus" class="mr-1 size-4" />
          新建任务
        </Button>
      </template>

      <template #progress="{ row }">
        <Progress
          :percent="row.progress"
          size="small"
          :status="
            row.status === 3
              ? 'exception'
              : row.status === 2
                ? 'success'
                : 'active'
          "
        />
      </template>

      <template #operation="{ row }">
        <Space>
          <Button
            :disabled="row.status !== 2"
            size="small"
            type="link"
            @click="onDownloadResult(row)"
          >
            <IconifyIcon icon="lucide:download" class="size-4" />
            下载结果
          </Button>
          <Popconfirm title="确定要删除该任务吗？" @confirm="onDelete(row)">
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
