import type { AppRouteRecordRaw } from "#src/router/types";

/**
 * 动态生成路由 - 前端方式
 * 权限系统已禁用：直接返回全部路由
 */
export function generateRoutesByFrontend(
	routes: AppRouteRecordRaw[],
	_roles: string[],
) {
	return routes;
}
