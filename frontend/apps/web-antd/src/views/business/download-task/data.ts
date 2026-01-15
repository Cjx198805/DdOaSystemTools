import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridOptions } from '#/adapter/vxe-table';

export function useFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'name',
      label: '任务名称',
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'url',
      label: '下载地址',
      rules: 'required',
    },
  ];
}

export function useGridFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'name',
      label: '任务名称',
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
      title: '任务名称',
      width: 200,
    },
    {
      field: 'url',
      title: '下载地址',
      minWidth: 300,
    },
    {
      field: 'progress',
      title: '进度',
      width: 200,
      slots: { default: 'progress' },
    },
    {
      cellRender: {
        name: 'CellTag',
        props: (row: any) => {
          const map: any = { 0: '等待中', 1: '进行中', 2: '已完成', 3: '失败' };
          const colorMap: any = {
            0: 'default',
            1: 'processing',
            2: 'success',
            3: 'error',
          };
          return {
            color: colorMap[row.status] || 'default',
            text: map[row.status] || '未知',
          };
        },
      },
      field: 'status',
      title: '状态',
      width: 100,
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
      width: 150,
    },
  ];
}
