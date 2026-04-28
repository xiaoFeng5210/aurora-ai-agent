import type { AppRouteRecordRaw } from "#src/router/types";

import { lazy } from "react";

import ContainerLayout from "#src/layout/container-layout";
import { about } from "#src/router/extra-info";

const About = lazy(() => import("#src/pages/about"));

const routes: AppRouteRecordRaw[] = [
	{
		path: "/about",
		Component: ContainerLayout,
		handle: {
			icon: "CopyrightOutlined",
			title: "common.menu.about",
			order: about,
		},
		children: [
			{
				index: true,
				Component: About,
				handle: {
					icon: "CopyrightOutlined",
					title: "common.menu.about",
				},
			},
		],
	},
];

export default routes;
