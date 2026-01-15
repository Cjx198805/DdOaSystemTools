import { requestClient } from '#/api/request';

export namespace SystemCompanyApi {
  export interface Company {
    id: number;
    name: string;
    code: string;
    type: number;
    status: number;
    created_at: string;
  }

  export interface ListParams {
    page?: number;
    page_size?: number;
    name?: string;
    code?: string;
  }

  export interface ListResult {
    list: Company[];
    total: number;
  }
}

export async function getSystemCompanyList(
  params: SystemCompanyApi.ListParams,
) {
  return requestClient.get<SystemCompanyApi.ListResult>('/company', { params });
}

export async function createSystemCompany(
  data: Partial<SystemCompanyApi.Company>,
) {
  return requestClient.post('/company', data);
}

export async function updateSystemCompany(
  id: number,
  data: Partial<SystemCompanyApi.Company>,
) {
  return requestClient.put(`/company/${id}`, data);
}

export async function deleteSystemCompany(id: number) {
  return requestClient.delete(`/company/${id}`);
}
