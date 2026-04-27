import type { AppRouteRecordRaw, RouteFileModule } from "#src/router/types";

import { loginPath } from "#src/router/extra-info";
import { ascending } from "#src/router/utils/ascending";
import { mergeRouteModules } from "#src/router/utils/merge-route-modules";
import { traverseTreeValues } from "#src/utils/tree";
import { coreRoutes } from "./core";

// 外部路由文件
export const externalRouteFiles: RouteFileModule = import.meta.glob("./external/**/*.ts", { eager: true });

/**
 * 前端业务路由文件。
 * 原模板把 modules 目录命名为“后端动态路由文件”，这里改为纯前端路由来源，
 * 不再依赖后端返回路由表。
 */
export const frontendRouteFiles: RouteFileModule = import.meta.glob("./modules/**/*.ts", { eager: true });

/**
 * 外部路由 1. 不进行权限校验， 2. 不会触发请求，例如用户信息接口
 * @example "privacy-policy", "terms-of-service" 等
 */
export const externalRoutes: AppRouteRecordRaw[] = mergeRouteModules(externalRouteFiles);

/** 前端业务路由 */
export const frontendRoutes: AppRouteRecordRaw[] = mergeRouteModules(frontendRouteFiles);

/**
 * 基本路由列表，由核心路由、外部路由组成，会一直存在系统中
 */
const baseRoutes = ascending([
	...coreRoutes,
	...externalRoutes,
]);

/** 权限路由列表，全部来自前端本地路由文件 */
const accessRoutes = [
	...frontendRoutes,
];

/**
 * 路由白名单 1. 不进行权限校验， 2. 不会触发请求，例如用户信息接口
 * @example "privacy-policy", "terms-of-service" 等
 */
const whiteRouteNames = [
	loginPath,
	...traverseTreeValues(externalRoutes, route => route.path),
];

export {
	accessRoutes,
	baseRoutes,
	whiteRouteNames,
};
