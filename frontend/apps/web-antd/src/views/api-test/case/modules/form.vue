<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useVbenDrawer } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { createTestCase, updateTestCase, type ApiTestApi } from '#/api/business/api_test';
import { useFormSchema } from '../data';

const emits = defineEmits(['success']);

const formData = ref<ApiTestApi.TestCase>();
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
      if (id.value) {
        await updateTestCase(id.value, values);
      } else {
        await createTestCase(values);
      }
      emits('success');
      drawerApi.close();
    } catch (error) {
    } finally {
      drawerApi.unlock();
    }
  },

  async onOpenChange(isOpen) {
    if (isOpen) {
      const data = drawerApi.getData<ApiTestApi.TestCase>();
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
  return id.value ? '编辑测试用例' : '新增测试用例';
});
</script>

<template>
  <Drawer :title="getDrawerTitle">
    <Form />
  </Drawer>
</template>
