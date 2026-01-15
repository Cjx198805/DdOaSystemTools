import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridOptions } from '#/adapter/vxe-table';

export function useFormSchema(): VbenFormSchema[] {
    return [
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
            component: 'Input',
            fieldName: 'label',
            label: '显示文本',
            rules: 'required',
        },
        {
            component: 'Input',
            fieldName: 'value',
            label: '值',
            rules: 'required',
        },
        {
            component: 'InputNumber',
            fieldName: 'sort',
            label: '排序',
            defaultValue: 0,
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
            fieldName: 'field',
            label: '字段',
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
            field: 'module',
            title: '模块',
            width: 120,
        },
        {
            field: 'field',
            title: '字段',
            width: 120,
        },
        {
            field: 'label',
            title: '显示文本',
            width: 150,
        },
        {
            field: 'value',
            title: '值',
            width: 100,
        },
        {
            field: 'sort',
            title: '排序',
            width: 80,
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
