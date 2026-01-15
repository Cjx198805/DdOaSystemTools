import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridOptions } from '#/adapter/vxe-table';

export function useFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'role_id',
      label: '角色ID',
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'module',
      label: '模块',
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'field',
      label: '字段',
      rules: 'required',
    },
    {
      component: 'Checkbox',
      fieldName: 'readable',
      label: '可读',
      defaultValue: 1,
      renderComponentContent: '可读',
    },
    {
      component: 'Checkbox',
      fieldName: 'editable',
      label: '可写',
      defaultValue: 1,
      renderComponentContent: '可写',
    },
    {
      component: 'Checkbox',
      fieldName: 'special_edit',
      label: '特殊编辑',
      defaultValue: 0,
      renderComponentContent: '特殊编辑',
    },
  ];
}

export function useGridFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'module',
      label: '模块',
    },
    {
      component: 'Input',
      fieldName: 'role_id',
      label: '角色ID',
    },
  ];
}

export function useColumns(): VxeTableGridOptions['columns'] {
  return [
    {
      field: 'id',
      title: 'ID',
      width: 80,
    },
    {
      field: 'role_id',
      title: '角色ID',
      width: 100,
    },
    {
      field: 'module',
      title: '模块名',
      width: 150,
    },
    {
      field: 'field',
      title: '显示名/代码',
      width: 150,
    },
    {
      cellRender: {
        name: 'CellTag',
        props: (row: any) => ({
          color: row.readable === 1 ? 'success' : 'error',
          text: row.readable === 1 ? '是' : '否',
        }),
      },
      field: 'readable',
      title: '可读',
      width: 80,
    },
    {
      cellRender: {
        name: 'CellTag',
        props: (row: any) => ({
          color: row.editable === 1 ? 'success' : 'error',
          text: row.editable === 1 ? '是' : '否',
        }),
      },
      field: 'editable',
      title: '可写',
      width: 80,
    },
    {
      cellRender: {
        name: 'CellTag',
        props: (row: any) => ({
          color: row.special_edit === 1 ? 'warning' : 'default',
          text: row.special_edit === 1 ? '是' : '否',
        }),
      },
      field: 'special_edit',
      title: '特殊编辑',
      width: 100,
    },
    {
      align: 'center',
      field: 'operation',
      fixed: 'right',
      slots: { default: 'operation' },
      title: '操作',
      width: 150,
    },
  ];
}
