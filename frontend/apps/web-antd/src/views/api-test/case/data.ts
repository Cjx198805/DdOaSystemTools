import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridOptions } from '#/adapter/vxe-table';

export function useFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'name',
      label: '用例名称',
      rules: 'required',
    },
    {
      component: 'Select',
      componentProps: {
        options: [
          { label: 'GET', value: 'GET' },
          { label: 'POST', value: 'POST' },
          { label: 'PUT', value: 'PUT' },
          { label: 'DELETE', value: 'DELETE' },
        ],
      },
      fieldName: 'method',
      label: '方法',
      defaultValue: 'GET',
    },
    {
      component: 'Input',
      fieldName: 'url',
      label: 'URL',
      rules: 'required',
    },
    {
      component: 'Textarea',
      fieldName: 'headers',
      label: 'Headers (JSON)',
    },
    {
      component: 'Textarea',
      fieldName: 'body',
      label: 'Body (JSON)',
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
      field: 'name',
      title: '用例名称',
      width: 200,
    },
    {
      cellRender: {
        name: 'CellTag',
        props: (row: any) => {
          const colors: any = {
            GET: 'green',
            POST: 'blue',
            PUT: 'orange',
            DELETE: 'red',
          };
          return { color: colors[row.method] || 'default', text: row.method };
        },
      },
      field: 'method',
      title: '方法',
      width: 100,
    },
    {
      field: 'url',
      title: 'URL',
      minWidth: 300,
    },
    {
      field: 'created_at',
      title: '创建时间',
      width: 180,
    },
    {
      align: 'center',
      field: 'operation',
      fixed: 'right',
      slots: { default: 'operation' },
      title: '操作',
      width: 250,
    },
  ];
}
