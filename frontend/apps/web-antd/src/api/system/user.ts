import { requestClient } from '#/api/request';

export namespace SystemUserApi {
  export interface User {
    id: number;
    username: string;
    nickname: string;
    email: string;
    phone: string;
    status: number;
    company_id: number;
    created_at: string;
  }

  export interface ListParams {
    page?: number;
    page_size?: number;
    username?: string;
    nickname?: string;
    status?: number;
    company_id?: number;
  }

  export interface ListResult {
    list: User[];
    total: number;
  }
}

/**
 * 获取系统用户列表
 */
export async function getSystemUserList(params: SystemUserApi.ListParams) {
  return requestClient.get<SystemUserApi.ListResult>('/user', { params });
}

/**
 * 创建用户
 */
export async function createSystemUser(data: Partial<SystemUserApi.User>) {
  return requestClient.post('/user', data);
}

/**
 * 更新用户
 */
export async function updateSystemUser(
  id: number,
  data: Partial<SystemUserApi.User>,
) {
  return requestClient.put(`/user/${id}`, data);
}

/**
 * 删除用户
 */
export async function deleteSystemUser(id: number) {
  return requestClient.delete(`/user/${id}`);
}

/**
 * 重置密码
 */
export async function resetUserPassword(id: number, password: string) {
  return requestClient.put(`/user/${id}/reset-password`, { password });
}
