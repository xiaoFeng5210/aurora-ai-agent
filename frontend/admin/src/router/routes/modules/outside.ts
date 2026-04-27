import type { AppRouteRecordRaw } from "#src/router/types";

import { Outlet } from "react-router";

import { Iframe } from "#src/components/iframe";
import { outside } from "#src/router/extra-info";

const routes: AppRouteRecordRaw[] = [
	{
		path: "/outside",
		Component: Outlet,
		handle: {
			icon: "OutsidePageIcon",
			title: "common.menu.outside",
			order: outside,
		},
		children: [
			{
				path: "/outside/embedded",
				Component: Outlet,
				handle: {
					icon: "EmbeddedIcon",
					title: "common.menu.embedded",
				},
				children: [
					{
						path: "/outside/embedded/ant-design",
						Component: Iframe,
						handle: {
							icon: "AntDesignOutlined",
							title: "common.menu.antd",
							iframeLink: "https://ant.design/",
						},
					},
					{
						path: "/outside/embedded/project-docs",
						Component: Iframe,
						handle: {
							icon: "ContainerOutlined",
							title: "common.menu.projectDocs",
							iframeLink: "https://condorheroblog.github.io/react-antd-admin/docs/",
						},
					},
				],
			},
			{
				path: "/outside/external-link",
				Component: Outlet,
				handle: {
					icon: "ExternalIcon",
					title: "common.menu.externalLink",
				},
				children: [
					{
						path: "/outside/external-link/react-docs",
						handle: {
							icon: "RiReactjsLine",
							title: "common.menu.reactDocs",
							externalLink: "https://react.dev/",
						},
					},
				],
			},
		],
	},
];

export default routes;
