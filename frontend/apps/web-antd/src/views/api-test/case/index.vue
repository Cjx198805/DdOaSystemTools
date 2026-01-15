<script lang="ts" setup>
import { Page, useVbenDrawer } from '@vben/common-ui';
import { IconifyIcon } from '@vben/icons';
import { Button, message, Popconfirm, Space } from 'ant-design-vue';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  getTestCaseList,
  deleteTestCase,
  runTestCase,
  type ApiTestApi,
} from '#/api/business/api_test';
import { useColumns } from './data';
import CaseForm from './modules/form.vue';

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  connectedComponent: CaseForm,
});

const [Grid, gridApi] = useVbenVxeGrid({
  gridOptions: {
    columns: useColumns(),
    height: 'auto',
    proxyConfig: {
      ajax: {
        query: async ({ page }) => {
          const res = await getTestCaseList({
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

function onRefresh() {
  gridApi.query();
}

function onCreate() {
  formDrawerApi.open();
}

function onEdit(row: ApiTestApi.TestCase) {
  formDrawerApi.setData(row).open();
}

async function onDelete(row: ApiTestApi.TestCase) {
  try {
    await deleteTestCase(row.id);
    message.success('删除成功');
    onRefresh();
  } catch (error) {}
}

async function onRun(row: ApiTestApi.TestCase) {
  const hide = message.loading('测试执行中...', 0);
  try {
    await runTestCase(row.id);
    message.success('测试运行完成，请前往历史查看结果');
  } catch (error) {
  } finally {
    hide();
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
          新增用例
        </Button>
      </template>

      <template #operation="{ row }">
        <Space>
          <Button size="small" type="link" @click="onRun(row)">
            <IconifyIcon icon="lucide:play-circle" class="size-4" />
            运行
          </Button>
          <Button size="small" type="link" @click="onEdit(row)">
            <IconifyIcon icon="lucide:notebook-pen" class="size-4" />
            编辑
          </Button>
          <Popconfirm title="确定要删除该测试用例吗？" @confirm="onDelete(row)">
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
