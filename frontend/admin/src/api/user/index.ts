import type { AppRouteRecordRaw } from "#src/router/types";
import type {
	AuthType,
	CreateUserPayload,
	DocumentRecord,
	LoginInfo,
	QueryUserDocumentPayload,
	QueryUserPayload,
	UserInfoType,
	UserRecord,
} from "./types";

import { request } from "#src/utils/request";
import { REFRESH_TOKEN_PATH } from "#src/utils/request/constants";

export * from "./types";

type GoApiResponse<T> = {
	code: number
	message: string
	data?: T
	result?: T
};

function getResponseData<T>(response: GoApiResponse<T>): T {
	return (response.data ?? response.result) as T;
}

export function fetchLogin(data: LoginInfo) {
	return request
		.post("login", { json: data })
		.json<GoApiResponse<unknown>>()
		.then((): ApiResponse<AuthType> => ({
			code: 0,
			message: "success",
			success: true,
			result: {
				token: "cookie-session",
				refreshToken: "",
			},
		}));
}

export function fetchLogout() {
	return request.post("logout").json();
}

export function fetchAsyncRoutes() {
	return Promise.resolve({
		code: 0,
		message: "success",
		success: true,
		result: [] as AppRouteRecordRaw[],
	});
}

export function fetchUserInfo() {
	return request
		.get("users/me")
		.json<GoApiResponse<UserRecord>>()
		.then((response): ApiResponse<UserInfoType> => {
			const user = getResponseData(response);
			return {
				code: 0,
				message: response.message,
				success: response.code === 0,
				result: {
					id: String(user.id),
					avatar: "",
					username: user.username,
					email: user.email,
					phoneNumber: user.phone,
					description: user.user_prompt,
					roles: ["admin"],
				},
			};
		});
}

export interface RefreshTokenResult {
	token: string
	refreshToken: string
}

export function fetchRefreshToken(data: { readonly refreshToken: string }) {
	return request.post(REFRESH_TOKEN_PATH, { json: data }).json<ApiResponse<RefreshTokenResult>>();
}

export function fetchUsers(data: QueryUserPayload = {}) {
	return request
		.post("users/query", { json: data })
		.json<GoApiResponse<UserRecord[]>>()
		.then(response => getResponseData(response) ?? []);
}

export function createUser(data: CreateUserPayload) {
	return request
		.post("register", { json: data })
		.json<GoApiResponse<unknown>>();
}

export function fetchUserDocuments(data: QueryUserDocumentPayload) {
	return request
		.post("documents/query-by-user", { json: data })
		.json<GoApiResponse<DocumentRecord[]>>()
		.then(response => getResponseData(response) ?? []);
}
