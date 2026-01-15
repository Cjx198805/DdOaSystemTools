import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridOptions } from '#/adapter/vxe-table';

export function useFormSchema(): VbenFormSchema[] {
    return [
        {
            component: 'Input',
            fieldName: 'name',
            label: '公司名称',
            rules: 'required',
        },
        {
            component: 'Input',
            fieldName: 'code',
            label: '公司编码',
            rules: 'required',
        },
        {
            component: 'Select',
            componentProps: {
                options: [
                    { label: '集团总部', value: 1 },
                    { label: '分子公司', value: 2 },
                ],
            },
            fieldName: 'type',
            label: '公司类型',
            defaultValue: 1,
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

export function useGridFormSchema(): VbenFormSchema[] {
    return [
        {
            component: 'Input',
            fieldName: 'name',
            label: '公司名称',
        },
        {
            component: 'Input',
            fieldName: 'code',
            label: '公司编码',
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
            title: '公司名称',
            width: 250,
        },
        {
            field: 'code',
            title: '公司编码',
            width: 150,
        },
        {
            cellRender: {
                name: 'CellTag',
                props: (row: any) => ({
                    color: row.type === 1 ? 'purple' : 'cyan',
                    text: row.type === 1 ? '集团总部' : '分子公司',
                }),
            },
            field: 'type',
            title: '类型',
            width: 120,
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
