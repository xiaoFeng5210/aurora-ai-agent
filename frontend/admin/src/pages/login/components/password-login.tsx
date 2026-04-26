import type { LoginInfo } from "#src/api/user";

import { useAuthStore } from "#src/store/auth";

import {
	Button,
	Form,
	Input,
	message,
	Space,
} from "antd";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useSearchParams } from "react-router";

const FORM_INITIAL_VALUES: LoginInfo = {
	email: "admin@example.com",
	password: "123456",
};

export function PasswordLogin() {
	const [loading, setLoading] = useState(false);
	const [passwordLoginForm] = Form.useForm();
	const { t } = useTranslation();
	const [messageLoadingApi, contextLoadingHolder] = message.useMessage();
	const [searchParams] = useSearchParams();
	const navigate = useNavigate();
	const login = useAuthStore(state => state.login);

	const handleFinish = async (values: LoginInfo) => {
		setLoading(true);
		messageLoadingApi?.loading(t("authority.loginInProgress"), 0);

		login(values).then(() => {
			messageLoadingApi?.destroy();
			window.$message?.success(t("authority.loginSuccess"));
			const redirect = searchParams.get("redirect");
			if (redirect) {
				navigate(`/${redirect.slice(1)}`);
			}
			else {
				navigate(import.meta.env.VITE_BASE_HOME_PATH);
			}
		}).finally(() => {
			messageLoadingApi?.destroy();
			// Prevent multiple requests from being made by clicking the login button
			setTimeout(() => {
				window.$message?.destroy();
				setLoading(false);
			}, 1000);
		});
	};

	return (
		<>
			{contextLoadingHolder}
			<Space orientation="vertical">
				<h2 className="text-colorText mb-3 text-3xl font-bold leading-9 tracking-tight lg:text-4xl">
					{t("authority.welcomeBack")}
					&nbsp;
					👏
				</h2>
				<p className="lg:text-base text-sm text-colorTextSecondary">
					{t("authority.loginDescription")}
				</p>
			</Space>

			<Form
				name="passwordLoginForm"
				form={passwordLoginForm}
				layout="vertical"
				initialValues={FORM_INITIAL_VALUES}
				onFinish={handleFinish}
			>
				<Form.Item
					label="邮箱"
					name="email"
					rules={[
						{ required: true, message: "请输入邮箱" },
						{ type: "email", message: "邮箱格式不正确" },
					]}
				>
					<Input placeholder="admin@example.com" />
				</Form.Item>

				<Form.Item
					label={t("authority.password")}
					name="password"
					rules={[{ required: true, message: "请输入密码" }]}
				>
					<Input.Password placeholder={t("form.password.required")} />
				</Form.Item>

				<Form.Item>
					<div className="flex justify-between mb-5 -mt-1 text-sm">
						<span className="text-colorTextSecondary">本地后台默认账号即可</span>
					</div>
					<Button block type="primary" htmlType="submit" loading={loading}>
						{t("authority.login")}
					</Button>
				</Form.Item>

			</Form>
		</>
	);
}
