import type { AppRouteRecordRaw } from "#src/router/types";
import ContainerLayout from "#src/layout/container-layout";
import { system } from "#src/router/extra-info";

import { lazy } from "react";

const User = lazy(() => import("#src/pages/system/user"));

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
		],
	},
];

export default routes;
