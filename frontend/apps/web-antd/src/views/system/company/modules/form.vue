<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useVbenDrawer } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { createSystemCompany, updateSystemCompany, type SystemCompanyApi } from '#/api/system/company';
import { useFormSchema } from '../data';

const emits = defineEmits(['success']);

const formData = ref<SystemCompanyApi.Company>();
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
        await updateSystemCompany(id.value, values);
      } else {
        await createSystemCompany(values);
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
      const data = drawerApi.getData<SystemCompanyApi.Company>();
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
  return id.value ? '编辑公司' : '新增公司';
});
</script>

<template>
  <Drawer :title="getDrawerTitle">
    <Form />
  </Drawer>
</template>
