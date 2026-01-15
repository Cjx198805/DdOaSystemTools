<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useVbenDrawer } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import {
  createFieldPermission,
  updateFieldPermission,
  type SystemFieldPermissionApi,
} from '#/api/system/field_permission';
import { useFormSchema } from '../data';

const emits = defineEmits(['success']);

const formData = ref<SystemFieldPermissionApi.FieldPermission>();
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
        await updateFieldPermission(id.value, values);
      } else {
        await createFieldPermission(values);
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
      const data =
        drawerApi.getData<SystemFieldPermissionApi.FieldPermission>();
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
  return id.value ? '编辑字段权限' : '新增字段权限';
});
</script>

<template>
  <Drawer :title="getDrawerTitle">
    <Form />
  </Drawer>
</template>
