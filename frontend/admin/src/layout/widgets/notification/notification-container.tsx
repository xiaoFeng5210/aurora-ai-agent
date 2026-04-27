import type { ButtonProps } from "antd";
import type { NotificationItem } from "./types";

import { fetchNotifications } from "#src/api/notifications";

import { useEffect, useState } from "react";
import { NotificationPopup } from "./index";

export function NotificationContainer({ ...restProps }: ButtonProps) {
	const [notifications, setNotifications] = useState<NotificationItem[]>([]);

	useEffect(() => {
		let ignore = false;

		fetchNotifications()
			.then((res) => {
				if (ignore) {
					return;
				}
				setNotifications(
					Array.from({ length: 20 }).flatMap(() => res.result ?? []),
				);
			})
			.catch(() => {
				if (!ignore) {
					setNotifications([]);
				}
			});

		return () => {
			ignore = true;
		};
	}, []);

	return (
		<NotificationPopup
			notifications={notifications}
			{...restProps}
		/>
	);
}
