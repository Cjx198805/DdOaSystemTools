import type { VxeTableGridOptions } from '#/adapter/vxe-table';

export function useColumns(): VxeTableGridOptions['columns'] {
  return [
    {
      field: 'id',
      title: 'ID',
      width: 80,
    },
    {
      field: 'case_name',
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
      minWidth: 200,
    },
    {
      cellRender: {
        name: 'CellTag',
        props: (row: any) => ({
          color:
            row.status_code >= 200 && row.status_code < 300
              ? 'success'
              : 'error',
          text: row.status_code,
        }),
      },
      field: 'status_code',
      title: '状态码',
      width: 100,
    },
    {
      field: 'duration',
      title: '耗时 (ms)',
      width: 120,
    },
    {
      field: 'created_at',
      title: '测试时间',
      width: 180,
    },
    {
      align: 'center',
      field: 'operation',
      fixed: 'right',
      slots: { default: 'operation' },
      title: '内容',
      width: 100,
    },
  ];
}
