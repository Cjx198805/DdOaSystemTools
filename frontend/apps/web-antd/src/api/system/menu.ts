import { requestClient } from '#/api/request';

export namespace SystemMenuApi {
  export interface Menu {
    id: number;
    parent_id: number;
    name: string;
    path: string;
    component: string;
    icon: string;
    sort: number;
    type: number;
    status: number;
    children?: Menu[];
  }
}

export async function getSystemMenuList() {
  return requestClient.get<SystemMenuApi.Menu[]>('/menu/all'); // Backend GetVbenTree or GetAll
}

export async function getSystemMenuTree() {
  return requestClient.get<SystemMenuApi.Menu[]>('/menu/tree');
}

export async function createSystemMenu(data: Partial<SystemMenuApi.Menu>) {
  return requestClient.post('/menu', data);
}

export async function updateSystemMenu(
  id: number,
  data: Partial<SystemMenuApi.Menu>,
) {
  return requestClient.put(`/menu/${id}`, data);
}

export async function deleteSystemMenu(id: number) {
  return requestClient.delete(`/menu/${id}`);
}
