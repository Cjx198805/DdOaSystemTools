<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useVbenDrawer } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { createSystemMenu, updateSystemMenu, type SystemMenuApi } from '#/api/system/menu';
import { useFormSchema } from '../data';

const emits = defineEmits(['success']);

const formData = ref<SystemMenuApi.Menu>();
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
        await updateSystemMenu(id.value, values);
      } else {
        await createSystemMenu(values);
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
      const data = drawerApi.getData<SystemMenuApi.Menu>();
      formApi.resetForm();

      if (data && data.id) {
        formData.value = data;
        id.value = data.id;
        formApi.setValues(data);
      } else {
        formData.value = undefined;
        id.value = undefined;
        // 如果有 parent_id 传进来，可以设置默认值
        if (data && data.parent_id) {
           formApi.setValues({ parent_id: data.parent_id });
        }
      }
    }
  },
});

const getDrawerTitle = computed(() => {
  return id.value ? '编辑菜单' : '新增菜单';
});
</script>

<template>
  <Drawer :title="getDrawerTitle">
    <Form />
  </Drawer>
</template>
