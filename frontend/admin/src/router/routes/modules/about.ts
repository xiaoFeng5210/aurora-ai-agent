import type { AppRouteRecordRaw } from "#src/router/types";

import { lazy } from "react";

import { about } from "#src/router/extra-info";

const About = lazy(() => import("#src/pages/about"));

const routes: AppRouteRecordRaw[] = [
	{
		path: "/about",
		Component: About,
		handle: {
			icon: "CopyrightOutlined",
			title: "common.menu.about",
			order: about,
		},
	},
];

export default routes;
