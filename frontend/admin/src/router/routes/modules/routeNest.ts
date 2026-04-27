import type { AppRouteRecordRaw } from "#src/router/types";

import { lazy } from "react";
import { Outlet } from "react-router";

import { routeNest } from "#src/router/extra-info";

const Menu11 = lazy(() => import("#src/pages/route-nest/menu1/menu1-1"));
const Menu12 = lazy(() => import("#src/pages/route-nest/menu1/menu1-2"));
const Menu2 = lazy(() => import("#src/pages/route-nest/menu2"));

const routes: AppRouteRecordRaw[] = [
	{
		path: "/route-nest",
		Component: Outlet,
		handle: {
			order: routeNest,
			title: "common.menu.nestMenus",
			icon: "NodeExpandOutlined",
		},
		children: [
			{
				path: "/route-nest/menu1",
				Component: Outlet,
				handle: {
					title: "common.menu.menu1",
					icon: "SisternodeOutlined",
				},
				children: [
					{
						path: "/route-nest/menu1/menu1-1",
						Component: Menu11,
						handle: {
							title: "common.menu.menu1-1",
							icon: "SubnodeOutlined",
						},
					},
					{
						path: "/route-nest/menu1/menu1-2",
						Component: Menu12,
						handle: {
							title: "common.menu.menu1-2",
							icon: "SubnodeOutlined",
						},
					},
				],
			},
			{
				path: "/route-nest/menu2",
				Component: Menu2,
				handle: {
					title: "common.menu.menu2",
					icon: "SubnodeOutlined",
				},
			},
		],
	},
];

export default routes;
