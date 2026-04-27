import type { AppRouteRecordRaw } from "#src/router/types";

import { lazy } from "react";
import { Outlet } from "react-router";

import { personalCenter } from "#src/router/extra-info";

const MyProfile = lazy(() => import("#src/pages/personal-center/my-profile"));
const Settings = lazy(() => import("#src/pages/personal-center/settings"));

const routes: AppRouteRecordRaw[] = [
	{
		path: "/personal-center",
		Component: Outlet,
		handle: {
			order: personalCenter,
			title: "common.menu.personalCenter",
			icon: "RiAccountCircleLine",
		},
		children: [
			{
				path: "/personal-center/my-profile",
				Component: MyProfile,
				handle: {
					title: "common.menu.profile",
					icon: "ProfileCardIcon",
				},
			},
			{
				path: "/personal-center/settings",
				Component: Settings,
				handle: {
					title: "common.menu.settings",
					icon: "RiUserSettingsLine",
				},
			},
		],
	},
];

export default routes;
