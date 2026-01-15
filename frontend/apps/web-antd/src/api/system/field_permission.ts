import { requestClient } from '#/api/request';

export namespace SystemFieldPermissionApi {
    export interface FieldPermission {
        id: number;
        role_id: number;
        module: string;
        field: string;
        readable: number;
        editable: number;
        special_edit: number;
    }

    export interface ListParams {
        page?: number;
        page_size?: number;
        role_id?: number;
        module?: string;
    }

    export interface ListResult {
        list: FieldPermission[];
        total: number;
    }
}

export async function getFieldPermissionList(params: SystemFieldPermissionApi.ListParams) {
    return requestClient.get<SystemFieldPermissionApi.ListResult>('/field-permission', { params });
}

export async function createFieldPermission(data: Partial<SystemFieldPermissionApi.FieldPermission>) {
    return requestClient.post('/field-permission', data);
}

export async function updateFieldPermission(id: number, data: Partial<SystemFieldPermissionApi.FieldPermission>) {
    return requestClient.put(`/field-permission/${id}`, data);
}

export async function deleteFieldPermission(id: number) {
    return requestClient.delete(`/field-permission/${id}`);
}
