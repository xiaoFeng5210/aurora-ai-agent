export * from "./constants";

/**
 * @zh 权限判断
 * 权限系统已禁用：所有判断函数恒返回 true，按钮 / 路由都不再做权限拦截。
 * 等后端权限就绪后，把此 Hook 还原即可恢复鉴权能力。
 */
export function useAccess() {
	const hasAccessByCodes = (_permission?: string | Array<string>) => true;
	const hasAccessByRoles = (_roles?: string | Array<string>) => true;
	return { hasAccessByCodes, hasAccessByRoles };
}
