import type { AppRouteRecordRaw } from "#src/router/types";

import { lazy } from "react";

import ContainerLayout from "#src/layout/container-layout";
import { system } from "#src/router/extra-info";

const User = lazy(() => import("#src/pages/system/user"));
const Role = lazy(() => import("#src/pages/system/role"));
const Menu = lazy(() => import("#src/pages/system/menu"));
const Dept = lazy(() => import("#src/pages/system/dept"));
const BaiduNetworkdisk = lazy(() => import("#src/pages/system/baidu-networkdisk"));

const routes: AppRouteRecordRaw[] = [
	{
		path: "/system",
		Component: ContainerLayout,
		handle: {
			icon: "SettingOutlined",
			title: "common.menu.system",
			order: system,
		},
		children: [
			{
				path: "/system/user",
				Component: User,
				handle: {
					icon: "UserOutlined",
					title: "common.menu.user",
				},
			},
			{
				path: "/system/role",
				Component: Role,
				handle: {
					icon: "TeamOutlined",
					title: "common.menu.role",
				},
			},
			{
				path: "/system/menu",
				Component: Menu,
				handle: {
					icon: "MenuOutlined",
					title: "common.menu.menu",
				},
			},
			{
				path: "/system/dept",
				Component: Dept,
				handle: {
					icon: "ApartmentOutlined",
					title: "common.menu.dept",
				},
			},
			{
				path: "/system/baidu-networkdisk",
				Component: BaiduNetworkdisk,
				handle: {
					icon: "CloudOutlined",
					title: "common.menu.baiduNetworkdisk",
				},
			},
		],
	},
];

export default routes;
