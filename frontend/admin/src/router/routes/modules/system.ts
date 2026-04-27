import type { AppRouteRecordRaw } from "#src/router/types";

import { lazy } from "react";

import ContainerLayout from "#src/layout/container-layout";
import { system } from "#src/router/extra-info";

const User = lazy(() => import("#src/pages/system/user"));
const BaiduNetworkdisk = lazy(() => import("#src/pages/system/baidu-networkdisk"));

const routes: AppRouteRecordRaw[] = [
	{
		path: "/system",
		Component: ContainerLayout,
		handle: {
			icon: "SettingOutlined",
			title: "common.menu.system",
			order: system,
			roles: ["admin"],
		},
		children: [
			{
				path: "/system/user",
				Component: User,
				handle: {
					icon: "UserOutlined",
					title: "common.menu.user",
					roles: ["admin"],
				},
			},
			{
				path: "/system/baidu-networkdisk",
				Component: BaiduNetworkdisk,
				handle: {
					icon: "CloudOutlined",
					title: "common.menu.baiduNetworkdisk",
					roles: ["admin"],
				},
			},
		],
	},
];

export default routes;
