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

export interface BaiduTokenResponse {
	expires_in?: number
	access_token?: string
	refresh_token?: string
	error?: string
	error_description?: string
}

export function exchangeBaiduNetworkdiskToken(code: string) {
	return request
		.post("file/baidu_networkdisk/token", { json: { code } })
		.json<GoApiResponse<BaiduTokenResponse>>()
		.then(response => getResponseData(response));
}
