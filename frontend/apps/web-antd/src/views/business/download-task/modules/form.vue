<script lang="ts" setup>
import { useVbenDrawer } from '@vben/common-ui';

import { useVbenForm } from '#/adapter/form';
import { createDownloadTask } from '#/api/business/download';

import { useFormSchema } from '../data';

const emits = defineEmits(['success']);

const [Form, formApi] = useVbenForm({
  schema: useFormSchema(),
  showDefaultActions: false,
});

const [Drawer, drawerApi] = useVbenDrawer({
  async onConfirm() {
    const { valid } = await formApi.validate();
    if (!valid) return;
    const values = await formApi.getValues();
    drawerApi.lock();
    try {
      await createDownloadTask(values);
      emits('success');
      drawerApi.close();
    } catch {
    } finally {
      drawerApi.unlock();
    }
  },

  async onOpenChange(isOpen) {
    if (isOpen) {
      formApi.resetForm();
    }
  },
});
</script>

<template>
  <Drawer title="新增下载任务">
    <Form />
  </Drawer>
</template>
