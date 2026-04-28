import {
	ProForm,
	ProFormDigit,
	ProFormText,
	ProFormTextArea,
} from "@ant-design/pro-components";
import { Card, Form, Input, Space, Typography } from "antd";
import { BasicContent } from "#src/components/basic-content";

import { FormAvatarItem } from "#src/components/basic-form";
import { useUserStore } from "#src/store/user";

const { Paragraph, Title } = Typography;

export default function Profile() {
	const currentUser = useUserStore();
	const getAvatarURL = () => {
		if (currentUser) {
			if (currentUser.avatar) {
				return currentUser.avatar;
			}
			const url = "https://avatar.vercel.sh/blur.svg?text=2";
			return url;
		}
		return "";
	};

	const handleFinish = async () => {
		window.$message?.success("更新基本信息成功");
	};

	return (
		<BasicContent className="min-h-full">
			<Space direction="vertical" size={16} className="w-full max-w-3xl">
				<Card>
					<Title level={3} className="!mb-1">我的资料</Title>
					<Paragraph type="secondary" className="!mb-0">
						维护当前账号的基础信息和展示资料。
					</Paragraph>
				</Card>

				<Card>
					<ProForm
						layout="vertical"
						onFinish={handleFinish}
						initialValues={{
							...currentUser,
							avatar: getAvatarURL(),
						}}
						requiredMark
					>
						<Form.Item
							name="avatar"
							label="头像"
							rules={[
								{
									required: true,
									message: "请上传头像",
								},
							]}
						>
							<FormAvatarItem />
						</Form.Item>
						<ProFormText
							name="username"
							label="用户名"
							rules={[
								{
									required: true,
									message: "请输入您的用户名!",
								},
							]}
						/>
						<ProFormText
							name="email"
							label="邮箱"
							rules={[
								{
									required: true,
									message: "请输入您的邮箱!",
								},
							]}
						/>
						<ProFormDigit
							name="phoneNumber"
							label="联系电话"
							rules={[
								{
									required: true,
									message: "请输入您的联系电话!",
								},
							]}
						>
							<Input type="tel" allowClear />
						</ProFormDigit>
						<ProFormTextArea
							allowClear
							name="description"
							label="个人简介"
							placeholder="个人简介"
						/>
					</ProForm>
				</Card>
			</Space>
		</BasicContent>
	);
};
