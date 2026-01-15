import { requestClient } from '#/api/request';

export namespace BusinessDownloadApi {
    export interface Task {
        id: number;
        user_id: number;
        name: string;
        url: string;
        status: number;
        progress: number;
        result_url: string;
        error_msg: string;
        created_at: string;
    }

    export interface ListParams {
        page?: number;
        page_size?: number;
    }

    export interface ListResult {
        list: Task[];
        total: number;
    }
}

export async function getDownloadTaskList(params: BusinessDownloadApi.ListParams) {
    return requestClient.get<BusinessDownloadApi.ListResult>('/download-task', { params });
}

export async function createDownloadTask(data: Partial<BusinessDownloadApi.Task>) {
    return requestClient.post('/download-task', data);
}

export async function deleteDownloadTask(id: number) {
    return requestClient.delete(`/download-task/${id}`);
}

export async function getDownloadResult(taskId: number) {
    return requestClient.get(`/download-task/result/${taskId}`);
}
