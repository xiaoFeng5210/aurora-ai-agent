import type { UserInfoType } from "#src/api/user/types";

import { hideLoading } from "#src/plugins/hide-loading";
import { setupLoading } from "#src/plugins/loading";
import { exception404Path, exception500Path, loginPath } from "#src/router/extra-info";
import { accessRoutes, whiteRouteNames } from "#src/router/routes";
import { generateRoutesByFrontend } from "#src/router/utils/generate-routes-from-frontend";
import { useAccessStore } from "#src/store/access";
import { useAuthStore } from "#src/store/auth";
import { useUserStore } from "#src/store/user";
import { goLogin } from "#src/utils/request/go-login";

import { useEffect } from "react";
import { matchRoutes, Navigate, useLocation, useNavigate, useSearchParams } from "react-router";

import { removeDuplicateRoutes } from "./utils";

/**
 * @zh 路由白名单 1. 不进行权限校验， 2. 不会触发请求，例如用户信息接口
 * @en Routes whitelist 1. No permission verification, 2. Will not trigger requests, such as user information interface
 * @example "privacy-policy", "terms-of-service" and so on.
 */
const noLoginWhiteList = Array.from(whiteRouteNames).filter(item => item !== loginPath);

interface AuthGuardProps {
	children?: React.ReactNode
}

/**
 * @zh AuthGuard 组件，用于权限验证，代码的顺序很重要，不要随意调整
 * @en AuthGuard component, used for permission verification. The order of the code is important and should not be arbitrarily adjusted
 */
export function AuthGuard({ children }: AuthGuardProps) {
	const navigate = useNavigate();
	const [searchParams] = useSearchParams();
	const { pathname, search } = useLocation();
	const isLogin = useAuthStore(state => Boolean(state.token));
	const isAuthorized = useUserStore(state => Boolean(state.id));
	const getUserInfo = useUserStore(state => state.getUserInfo);
	const { setAccessStore, isAccessChecked, routeList } = useAccessStore();

	const isPathInNoLoginWhiteList = noLoginWhiteList.includes(pathname);

	/**
	 * @zh 异步获取用户信息，并使用前端本地路由生成菜单和路由表。
	 * @en Fetch user information and use frontend local routes to generate menus and routes.
	 */
	useEffect(() => {
		async function fetchUserInfoAndRoutes() {
			/**
			 * @zh 登录跳转，防止闪烁
			 * @en Login redirect, prevent flicker
			 */
			setupLoading();

			type UserInfoResult =
				| { status: "fulfilled", value: UserInfoType }
				| { status: "rejected", reason: any };

			const userInfoResult: UserInfoResult = await Promise.resolve(getUserInfo()).then(
				value => ({ status: "fulfilled", value }),
				reason => ({ status: "rejected", reason }),
			);

			const latestRoles: string[] = [];
			/**
			 * @zh 从用户接口中获取角色信息
			 * @en Fetch role information from the user interface
			 */
			if (userInfoResult.status === "fulfilled" && "roles" in userInfoResult.value) {
				latestRoles.push(...userInfoResult.value?.roles ?? []);
			}

			const frontendRoutes = generateRoutesByFrontend(accessRoutes, latestRoles);
			const uniqueRoutes = removeDuplicateRoutes(frontendRoutes);
			setAccessStore(uniqueRoutes);

			/**
			 * @zh 用户信息接口异常：401/403 视为登录失效，跳转登录页；其他跳转 500。
			 * @en Handle user info request failure: 401/403 treated as auth failure (go login); others → 500.
			 */
			if (userInfoResult.status === "rejected") {
				const status = userInfoResult.reason?.response?.status;
				if (status === 401 || status === 403) {
					goLogin();
					return;
				}
				return navigate(exception500Path);
			}

			/**
			 *
			 * @zh 开启动态路由条件下需要替换当前路由？
			 * 1. 浏览器导航进入动态路由地址，例如 /system/user
			 * 2. 动态路由未添加到路由，所以地址栏中依然是 /system/user 但匹配到的路由是 fallback (path = "*") 路由
			 * 3. 添加完动态路由后，使用 replace 替换当前路由，触发程序重新匹配到 /system/user 路由
			 *
			 * Refer：https://router.vuejs.org/guide/advanced/dynamic-routing#Adding-routes
			 *
			 * @en Under the condition of dynamic routing, do you need to replace the current route?
			 * 1. Browser navigation into a dynamic routing address, such as /system/user
			 * 2. The dynamic route is not added to the route, so the address bar is still /system/user but the matched route is the fallback (path = "*") route
			 * 3. After adding the dynamic route, use replace to replace the current route and trigger the program to match /system/user again
			 */
			navigate(`${pathname}${search}`, {
				replace: true,
				/**
				 * @zh 保证替换路由前不会显示 404 页面（登录页面，网速切换为 3G 会闪烁显示 404 页面）
				 * @en Ensure that the 404 page will not be displayed before replacing the route
				 */
				flushSync: true,
			});
		}
		/**
		 * @zh 只有在以下条件下才执行获取用户信息和路由的逻辑
		 * 1. 非路由白名单
		 * 2. 已登录
		 * 3. 未获取到用户信息和路由信息
		 *
		 * @en The logic of obtaining user information and routes is only executed under the following conditions
		 * 1. Not in the route whitelist
		 * 2. Logged in
		 * 3. Unable to obtain user information and route information
		 *
		 */
		if (!whiteRouteNames.includes(pathname) && isLogin && !isAuthorized) {
			fetchUserInfoAndRoutes();
		}
	}, [pathname, search, isLogin, isAuthorized, getUserInfo, setAccessStore, navigate]);

	/**
	 * @zh 路由白名单
	 * @en Route whitelist
	 * @see {noLoginWhiteList}
	 */
	if (isPathInNoLoginWhiteList) {
		hideLoading();
		return children;
	}

	/**
	 * @zh 未登录条件下的处理逻辑
	 * @en Processing logic under unlogged conditions
	 */
	/* --------------- Start ------------------ */
	if (!isLogin) {
		hideLoading();
		// 未登录且目标页不是登录页，则跳转到登录页
		if (pathname !== loginPath) {
			// pathname 长度大于 1，则携带当前路径跳转登录页，否则直接跳转登录页
			const redirectPath = pathname.length > 1 ? `${loginPath}?redirect=${pathname}${search}` : loginPath;
			return (
				<Navigate
					to={redirectPath}
					replace
				/>
			);
		}
		// 未登录且目标页是登录页，保留登录页
		else {
			return children;
		}
	}
	/* --------------- End ------------------ */

	/**
	 * @zh 登录条件下的处理逻辑
	 * @en Processing logic under logged conditions
	 */
	/* --------------- Start ------------------ */

	/**
	 * @zh 已登录条件下，匹配 login 路由，跳转到首页
	 * 放到用户信息前，因为 login 路由不会请求用户信息，所以放在前面判断
	 *
	 * @en Under logged conditions, match the login route and jump to the home page
	 * Put it before user information, because the login route will not request user information, so put it in front to judge
	 */
	if (pathname === loginPath) {
		/**
		 * @example login?redirect=/system/user
		 */
		const redirectPath = searchParams.get("redirect");
		if (redirectPath?.length && redirectPath !== pathname) {
			return (
				<Navigate
					to={redirectPath}
					replace
				/>
			);
		}
		return (
			<Navigate
				to={import.meta.env.VITE_BASE_HOME_PATH}
				replace
			/>
		);
	}

	/**
	 * @zh 等待获取用户信息
	 * @en  Waiting for user information to be obtained
	 */
	if (!isAuthorized) {
		return null;
	}
	/**
	 * @zh 等待获取路由信息
	 * @en Waiting for route information to be obtained
	 */
	if (!isAccessChecked) {
		return null;
	}

	/**
	 * @zh 隐藏加载动画
	 * @en Hide loading animation
	 */
	hideLoading();

	/**
	 * @zh 如果是根路由则跳转到首页（获取完用户信息之后跳转到默认首页，防止请求两次用户信息接口）
	 * @en If it is the root route, jump to the home page (jump to the default home page after obtaining user information to prevent requesting twice for user information interface)
	 * @zh pathname 返回的是相对 import.meta.env.BASE_URL 的路径，所以这里是相对于 BASE_URL 的根路由 "/"
	 * @en pathname returns the path relative to import.meta.env.BASE_URL, so here is the root route "/" relative to BASE_URL
	 */
	if (pathname === "/") {
		return (
			<Navigate
				to={import.meta.env.VITE_BASE_HOME_PATH}
				replace
			/>
		);
	}

	/* --------------- End ------------------ */

	const matches = matchRoutes(
		routeList,
		pathname,
	) ?? [];

	const hasChildren = matches[matches.length - 1]?.route?.children?.filter(item => !item.index)?.length;
	const isMatchedByRouteList = matches.length > 0 && matches[matches.length - 1]?.route?.path !== "*";

	/**
	 * @zh 如果当前路径没有匹配到任何业务路由，则优先跳转到首页，避免后端未返回路由时直接落到 404
	 * @en If the current path does not match any business route, redirect to the home page first to avoid landing on 404 when the backend does not return routes
	 */
	if (!isMatchedByRouteList) {
		return (
			<Navigate
				to={import.meta.env.VITE_BASE_HOME_PATH}
				replace
			/>
		);
	}

	/**
	 * @zh 如果当前路由有子路由，则跳转到 404 页面
	 */
	if (hasChildren && hasChildren > 0) {
		return (
			<Navigate
				to={exception404Path}
				replace
			/>
		);
	}

	return children;
}
/**
 * 验证路由跳转是否正确的步骤：
 * 1. 未登录情况下，输入 login 路由
 * 2. 未登录情况下，输入非 login 路由
 * 3. 已登录情况下，使用系统的退出登录，然后再次登录
 * 4. 任选一个非 home 页面，使用开发者工具清除 localStorage，刷新页面之后进行登录
 * 5. 已登录情况下，输入 login 路由
 * 6. 已登录情况下，输入非 login 路由
 * 7. 已登录情况下，输入 http://localhost:3333 跳转到 /home 路由，用户接口发送一次
 * 8. 已登录情况下，输入 http://localhost:3333/ 跳转到 /home 路由，用户接口发送一次
 * 9. 已登录情况下，输入 http://localhost:3333/home 跳转到 /home 路由，用户接口发送一次
 */
