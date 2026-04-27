import type { AppRouteRecordRaw } from "#src/router/types";

import { lazy } from "react";
import { Outlet } from "react-router";

import { access } from "#src/router/extra-info";

const AccessMode = lazy(() => import("#src/pages/access/access-mode"));
const PageControl = lazy(() => import("#src/pages/access/page-control"));
const ButtonControl = lazy(() => import("#src/pages/access/button-control"));
const AdminVisible = lazy(() => import("#src/pages/access/admin-visible"));
const CommonVisible = lazy(() => import("#src/pages/access/common-visible"));

const routes: AppRouteRecordRaw[] = [
	{
		path: "/access",
		Component: Outlet,
		handle: {
			icon: "SafetyOutlined",
			title: "common.menu.access",
			order: access,
		},
		children: [
			{
				path: "/access/access-mode",
				Component: AccessMode,
				handle: {
					icon: "CloudOutlined",
					title: "common.menu.accessMode",
				},
			},
			{
				path: "/access/page-control",
				Component: PageControl,
				handle: {
					icon: "FileTextOutlined",
					title: "common.menu.pageControl",
				},
			},
			{
				path: "/access/button-control",
				Component: ButtonControl,
				handle: {
					icon: "LockOutlined",
					title: "common.menu.buttonControl",
					permissions: [
						"permission:button:get",
						"permission:button:update",
						"permission:button:delete",
						"permission:button:add",
					],
				},
			},
			{
				path: "/access/admin-visible",
				Component: AdminVisible,
				handle: {
					icon: "EyeOutlined",
					title: "common.menu.adminVisible",
				},
			},
			{
				path: "/access/common-visible",
				Component: CommonVisible,
				handle: {
					icon: "EyeOutlined",
					title: "common.menu.commonVisible",
				},
			},
		],
	},
];

export default routes;
