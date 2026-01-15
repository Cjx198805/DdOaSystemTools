<script lang="ts" setup>
import type { ApiTestApi } from '#/api/business/api_test';

import { Page } from '@vben/common-ui';
import { IconifyIcon } from '@vben/icons';

import { Button, Modal } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getTestHistoryList } from '#/api/business/api_test';

import { useColumns } from './data';

const [Grid, gridApi] = useVbenVxeGrid({
  gridOptions: {
    columns: useColumns(),
    height: 'auto',
    proxyConfig: {
      ajax: {
        query: async ({ page }) => {
          const res = await getTestHistoryList({
            page: page.currentPage,
            page_size: page.pageSize,
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

function onViewResponse(row: ApiTestApi.TestHistory) {
  Modal.info({
    title: '响应正文',
    content: row.response_body,
    width: 800,
    okText: '确定',
  });
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #operation="{ row }">
        <Button size="small" type="link" @click="onViewResponse(row)">
          <IconifyIcon icon="lucide:eye" class="mr-1 size-4" />
          查看详情
        </Button>
      </template>
    </Grid>
  </Page>
</template>
