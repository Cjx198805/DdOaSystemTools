<script lang="ts" setup>
import { computed, ref, nextTick } from 'vue';
import { useVbenDrawer, Tree } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import {
  createSystemRole,
  updateSystemRole,
  assignRoleMenus,
  getRoleMenus,
  type SystemRoleApi,
} from '#/api/system/role';
import { getSystemMenuList } from '#/api/system/menu';
import { useFormSchema } from '../data';
import { Spin } from 'ant-design-vue';

const emits = defineEmits(['success']);

const formData = ref<SystemRoleApi.Role>();
const id = ref<number>();
const menuTreeData = ref<any[]>([]);
const selectedMenuIds = ref<number[]>([]);
const loadingMenus = ref(false);

const [Form, formApi] = useVbenForm({
  schema: useFormSchema(),
  showDefaultActions: false,
});

async function loadMenuTree() {
  loadingMenus.value = true;
  try {
    const res = await getSystemMenuList();
    menuTreeData.value = res;
  } finally {
    loadingMenus.value = false;
  }
}

const [Drawer, drawerApi] = useVbenDrawer({
  async onConfirm() {
    const { valid } = await formApi.validate();
    if (!valid) return;
    const values = await formApi.getValues();
    drawerApi.lock();
    try {
      let roleId = id.value;
      if (roleId) {
        await updateSystemRole(roleId, values);
      } else {
        const res = await createSystemRole(values);
        roleId = (res as any).id;
      }

      // 分配菜单
      if (roleId) {
        await assignRoleMenus(roleId, selectedMenuIds.value);
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
      const data = drawerApi.getData<SystemRoleApi.Role>();
      formApi.resetForm();
      selectedMenuIds.value = [];

      if (menuTreeData.value.length === 0) {
        await loadMenuTree();
      }

      if (data && data.id) {
        formData.value = data;
        id.value = data.id;
        formApi.setValues(data);

        // 加载已绑定的菜单
        const res = await getRoleMenus(data.id);
        const menus = (res as any) || [];
        selectedMenuIds.value = menus.map((m: any) => m.id);
      } else {
        formData.value = undefined;
        id.value = undefined;
      }
    }
  },
});

const getDrawerTitle = computed(() => {
  return id.value ? '编辑角色' : '新增角色';
});
</script>

<template>
  <Drawer :title="getDrawerTitle" width="600">
    <Form />
    <div class="mt-4 rounded border p-4">
      <div class="mb-2 font-bold">菜单权限分配</div>
      <Spin :spinning="loadingMenus">
        <Tree
          v-model:value="selectedMenuIds"
          :tree-data="menuTreeData"
          checkable
          multiple
          :field-names="{
            children: 'children',
            title: 'name',
            key: 'id',
          }"
          style="max-height: 400px; overflow: auto"
        />
      </Spin>
    </div>
  </Drawer>
</template>
