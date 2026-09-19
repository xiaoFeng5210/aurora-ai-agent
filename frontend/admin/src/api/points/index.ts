import { request } from "#src/utils/request";

interface GoApiResponse<T> {
	code: number
	message: string
	data?: T
	result?: T
}

function getResponseData<T>(response: GoApiResponse<T>): T {
	return (response.data ?? response.result) as T;
}

export interface PointsBalance {
	balance: number
}

export interface AdjustPointsPayload {
	amount: number
	remark: string
}

export function fetchUserPoints(userId: number) {
	return request
		.get(`users/${userId}/points`)
		.json<GoApiResponse<PointsBalance>>()
		.then(response => getResponseData(response));
}

export function addUserPoints(userId: number, data: AdjustPointsPayload) {
	return request
		.post(`users/${userId}/points/add`, { json: data })
		.json<GoApiResponse<PointsBalance>>()
		.then(response => getResponseData(response));
}

export function deductUserPoints(userId: number, data: AdjustPointsPayload) {
	return request
		.post(`users/${userId}/points/deduct`, { json: data })
		.json<GoApiResponse<PointsBalance>>()
		.then(response => getResponseData(response));
}
