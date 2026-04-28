import { defineFakeRoute } from "vite-plugin-fake-server/client";
import { about, home, personalCenter, system } from "#/src/router/extra-info";
import { resultSuccess } from "./utils";

/**
 * roles：页面级别权限，这里模拟二种 "admin"、"common"
 * admin：管理员角色
 * common：普通角色
 */

const systemManagementRouter = {
	path: "/system",
	handle: {
		icon: "SettingOutlined",
		title: "common.menu.system",
		order: system,
		// roles: ["admin"],
	},
	children: [
		{
			path: "/system/user",
			component: "/system/user/index.tsx",
			handle: {
				icon: "UserOutlined",
				title: "common.menu.user",
				// permissions: [
				// 	"permission:button:add",
				// 	"permission:button:update",
				// 	"permission:button:delete",
				// ],
			},
		},
	],
};

const homeRouter = {
	path: "/home",
	component: "/home/index.tsx",
	handle: {
		icon: "HomeOutlined",
		title: "common.menu.home",
		order: home,
	},
};

const aboutRouter = {
	path: "/about",
	component: "/about/index.tsx",
	handle: {
		icon: "CopyrightOutlined",
		title: "common.menu.about",
		order: about,
	},
};

const personalCenterRouter = {
	path: "/personal-center",
	handle: {
		order: personalCenter,
		title: "common.menu.personalCenter",
		icon: "RiAccountCircleLine",
	},
	children: [
		{
			path: "/personal-center/my-profile",
			handle: {
				title: "common.menu.profile",
				icon: "ProfileCardIcon",
			},
		},
		{
			path: "/personal-center/settings",
			handle: {
				title: "common.menu.settings",
				icon: "RiUserSettingsLine",
			},
		},
	],
};

export default defineFakeRoute([
	{
		url: "/get-async-routes",
		timeout: 1000,
		method: "get",
		response: () => {
			return resultSuccess(
				[
					homeRouter,
					aboutRouter,
					systemManagementRouter,
					personalCenterRouter,
				],
			);
		},
	},
]);
