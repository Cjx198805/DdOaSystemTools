import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridOptions } from '#/adapter/vxe-table';

export function useFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'name',
      label: '菜单名称',
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'path',
      label: '路径',
    },
    {
      component: 'Input',
      fieldName: 'component',
      label: '组件',
    },
    {
      component: 'Input',
      fieldName: 'icon',
      label: '图标',
    },
    {
      component: 'InputNumber',
      fieldName: 'sort',
      label: '排序',
      defaultValue: 0,
    },
    {
      component: 'RadioGroup',
      componentProps: {
        buttonStyle: 'solid',
        options: [
          { label: '菜单', value: 1 },
          { label: '按钮', value: 2 },
        ],
        optionType: 'button',
      },
      defaultValue: 1,
      fieldName: 'type',
      label: '类型',
    },
    {
      component: 'RadioGroup',
      componentProps: {
        buttonStyle: 'solid',
        options: [
          { label: '启用', value: 1 },
          { label: '禁用', value: 0 },
        ],
        optionType: 'button',
      },
      defaultValue: 1,
      fieldName: 'status',
      label: '状态',
    },
  ];
}

export function useColumns(): VxeTableGridOptions['columns'] {
  return [
    {
      field: 'name',
      title: '菜单名称',
      width: 250,
      treeNode: true,
      align: 'left',
    },
    {
      field: 'path',
      title: '路径',
      width: 200,
    },
    {
      field: 'component',
      title: '组件',
      width: 200,
    },
    {
      cellRender: {
        name: 'CellTag',
        props: (row: any) => ({
          color: row.type === 1 ? 'blue' : 'orange',
          text: row.type === 1 ? '菜单' : '按钮',
        }),
      },
      field: 'type',
      title: '类型',
      width: 100,
    },
    {
      cellRender: {
        name: 'CellTag',
        props: (row: any) => ({
          color: row.status === 1 ? 'success' : 'error',
          text: row.status === 1 ? '启用' : '禁用',
        }),
      },
      field: 'status',
      title: '状态',
      width: 100,
    },
    {
      align: 'center',
      field: 'operation',
      fixed: 'right',
      slots: { default: 'operation' },
      title: '操作',
      width: 200,
    },
  ];
}
