import type { AppRouteRecordRaw } from "#src/router/types";

import { lazy } from "react";

import { home } from "#src/router/extra-info";

const Home = lazy(() => import("#src/pages/home"));

const routes: AppRouteRecordRaw[] = [
	{
		path: "/home",
		Component: Home,
		handle: {
			icon: "HomeOutlined",
			title: "common.menu.home",
			order: home,
		},
	},
];

export default routes;
