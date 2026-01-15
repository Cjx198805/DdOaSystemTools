import { requestClient } from '#/api/request';

export namespace SystemDictApi {
  export interface Dict {
    id: number;
    module: string;
    field: string;
    label: string;
    value: string;
    sort: number;
    status: number;
  }

  export interface ListParams {
    page?: number;
    page_size?: number;
    module?: string;
    field?: string;
  }

  export interface ListResult {
    list: Dict[];
    total: number;
  }
}

export async function getSystemDictList(params: SystemDictApi.ListParams) {
  return requestClient.get<SystemDictApi.ListResult>('/data-dictionary', {
    params,
  });
}

export async function createSystemDict(data: Partial<SystemDictApi.Dict>) {
  return requestClient.post('/data-dictionary', data);
}

export async function updateSystemDict(
  id: number,
  data: Partial<SystemDictApi.Dict>,
) {
  return requestClient.put(`/data-dictionary/${id}`, data);
}

export async function deleteSystemDict(id: number) {
  return requestClient.delete(`/data-dictionary/${id}`);
}
