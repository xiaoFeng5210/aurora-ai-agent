import type { AppRouteRecordRaw } from "#src/router/types";

export interface AuthType {
	token: string
	refreshToken: string
}

export interface LoginInfo {
	email?: string
	password: string
	username?: string
}

export interface UserInfoType {
	id: string
	avatar: string
	username: string
	email: string
	phoneNumber: string
	description: string
	roles: Array<string>
	// 路由可以在此处动态添加
	menus?: AppRouteRecordRaw[]
}

export interface AuthListProps {
	label: string
	name: string
	auth: string[]
}

export interface UserRecord {
	id: number
	username: string
	email: string
	phone: string
	birthday: string
	user_prompt: string
	created_at: string
	updated_at: string
}

export interface CreateUserPayload {
	username: string
	password: string
	email: string
	phone?: string
	birthday?: string
	user_prompt?: string
}

export interface QueryUserPayload {
	username?: string
	email?: string
	phone?: string
	birthday?: string
	page?: number
	page_size?: number
}

export interface DocumentRecord {
	id: number
	user_id: number
	display_name: string
	file_name?: string | null
	created_at: string
	updated_at: string
}

export interface QueryUserDocumentPayload {
	user_id: number
	display_name?: string
	file_name?: string
	page?: number
	page_size?: number
}
