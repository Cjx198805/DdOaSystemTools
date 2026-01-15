import { requestClient } from '#/api/request';

export namespace SystemRoleApi {
  export interface Role {
    id: number;
    name: string;
    code: string;
    status: number;
    created_at: string;
  }

  export interface ListParams {
    page?: number;
    page_size?: number;
    name?: string;
    code?: string;
    status?: number;
  }

  export interface ListResult {
    list: Role[];
    total: number;
  }
}

export async function getSystemRoleList(params: SystemRoleApi.ListParams) {
  return requestClient.get<SystemRoleApi.ListResult>('/role', { params });
}

export async function createSystemRole(data: Partial<SystemRoleApi.Role>) {
  return requestClient.post('/role', data);
}

export async function updateSystemRole(
  id: number,
  data: Partial<SystemRoleApi.Role>,
) {
  return requestClient.put(`/role/${id}`, data);
}

export async function deleteSystemRole(id: number) {
  return requestClient.delete(`/role/${id}`);
}

export async function getRoleMenus(id: number) {
  return requestClient.get(`/role/${id}/menus`);
}

export async function assignRoleMenus(id: number, menu_ids: number[]) {
  return requestClient.put(`/role/${id}/assign-menus`, { menu_ids });
}
