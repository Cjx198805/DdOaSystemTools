<script lang="ts" setup>
import type { SystemUserApi } from '#/api/system/user';

import { computed, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';

import { useVbenForm } from '#/adapter/form';
import { createSystemUser, updateSystemUser } from '#/api/system/user';

import { useFormSchema } from '../data';

const emits = defineEmits(['success']);

const formData = ref<SystemUserApi.User>();
const id = ref<number>();

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
      await (id.value
        ? updateSystemUser(id.value, values)
        : createSystemUser(values));
      emits('success');
      drawerApi.close();
    } catch {
      // Error handled by interceptor
    } finally {
      drawerApi.unlock();
    }
  },

  async onOpenChange(isOpen) {
    if (isOpen) {
      const data = drawerApi.getData<SystemUserApi.User>();
      formApi.resetForm();

      if (data && data.id) {
        formData.value = data;
        id.value = data.id;
        formApi.setValues(data);
      } else {
        formData.value = undefined;
        id.value = undefined;
      }
    }
  },
});

const getDrawerTitle = computed(() => {
  return id.value ? '编辑用户' : '新增用户';
});
</script>

<template>
  <Drawer :title="getDrawerTitle">
    <Form />
  </Drawer>
</template>
