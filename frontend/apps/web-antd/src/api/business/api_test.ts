import { requestClient } from '#/api/request';

export namespace ApiTestApi {
  export interface TestCase {
    id: number;
    name: string;
    description: string;
    method: string;
    url: string;
    headers: string;
    body: string;
    params: string;
    created_at: string;
  }

  export interface TestHistory {
    id: number;
    case_id: number;
    case_name: string;
    url: string;
    method: string;
    status_code: number;
    response_body: string;
    duration: number;
    created_at: string;
  }

  export interface CaseListParams {
    page?: number;
    page_size?: number;
    name?: string;
  }

  export interface HistoryListParams {
    page?: number;
    page_size?: number;
    case_id?: number;
  }
}

export async function getTestCaseList(params: ApiTestApi.CaseListParams) {
  return requestClient.get('/api-test/case', { params });
}

export async function createTestCase(data: Partial<ApiTestApi.TestCase>) {
  return requestClient.post('/api-test/case', data);
}

export async function updateTestCase(
  id: number,
  data: Partial<ApiTestApi.TestCase>,
) {
  return requestClient.put(`/api-test/case/${id}`, data);
}

export async function deleteTestCase(id: number) {
  return requestClient.delete(`/api-test/case/${id}`);
}

export async function runTestCase(id: number) {
  return requestClient.post(`/api-test/case/${id}/run`);
}

export async function getTestHistoryList(params: ApiTestApi.HistoryListParams) {
  return requestClient.get('/api-test/history', { params });
}
