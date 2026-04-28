import type { AppRouteRecordRaw } from "#src/router/types";

import { lazy } from "react";

import ContainerLayout from "#src/layout/container-layout";
import { home } from "#src/router/extra-info";

const Home = lazy(() => import("#src/pages/home"));

const routes: AppRouteRecordRaw[] = [
	{
		path: "/home",
		Component: ContainerLayout,
		handle: {
			icon: "HomeOutlined",
			title: "common.menu.home",
			order: home,
		},
		children: [
			{
				index: true,
				Component: Home,
				handle: {
					icon: "HomeOutlined",
					title: "common.menu.home",
				},
			},
		],
	},
];

export default routes;
