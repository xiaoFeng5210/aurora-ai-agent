import type { BaiduTokenResponse } from "#src/api/baidu-networkdisk";

import {
	CloudOutlined,
	KeyOutlined,
	SaveOutlined,
} from "@ant-design/icons";
import {
	Alert,
	Button,
	Descriptions,
	Form,
	Input,
	message,
	Space,
	Typography,
} from "antd";
import { useState } from "react";

import { exchangeBaiduNetworkdiskToken } from "#src/api/baidu-networkdisk";
import { BasicContent } from "#src/components/basic-content";

const { Text, Title } = Typography;

interface TokenForm {
	code: string
}

function maskToken(value?: string) {
	if (!value) {
		return "-";
	}
	if (value.length <= 16) {
		return value;
	}
	return `${value.slice(0, 8)}...${value.slice(-8)}`;
}

export default function BaiduNetworkdisk() {
	const [form] = Form.useForm<TokenForm>();
	const [saving, setSaving] = useState(false);
	const [token, setToken] = useState<BaiduTokenResponse>();

	const handleSubmit = async () => {
		const values = await form.validateFields();
		setSaving(true);
		try {
			const nextToken = await exchangeBaiduNetworkdiskToken(values.code.trim());
			setToken(nextToken);
			message.success("百度网盘 token 已保存到 Redis");
			form.resetFields();
		}
		finally {
			setSaving(false);
		}
	};

	return (
		<BasicContent className="min-h-full">
			<div className="flex max-w-4xl flex-col gap-4">
				<div>
					<Title level={3} className="!mb-1">百度网盘 Token</Title>
					<Text type="secondary">提交授权 code，后台会换取 token 并写入 Redis。</Text>
				</div>

				<section className="rounded-md border border-colorBorder bg-colorBgContainer p-4">
					<div className="mb-4 flex items-center gap-2">
						<KeyOutlined />
						<Text strong>授权 code</Text>
					</div>
					<Form form={form} layout="vertical" onFinish={handleSubmit}>
						<Form.Item
							label="Code"
							name="code"
							rules={[
								{ required: true, message: "请输入百度网盘授权 code" },
								{ whitespace: true, message: "Code 不能为空" },
							]}
						>
							<Input.Password
								placeholder="粘贴百度网盘授权 code"
								autoComplete="off"
							/>
						</Form.Item>
						<Form.Item className="!mb-0">
							<Button type="primary" htmlType="submit" icon={<SaveOutlined />} loading={saving}>
								保存 token
							</Button>
						</Form.Item>
					</Form>
				</section>

				<section className="rounded-md border border-colorBorder bg-colorBgContainer p-4">
					<div className="mb-4 flex items-center gap-2">
						<CloudOutlined />
						<Text strong>保存结果</Text>
					</div>
					{token
						? (
							<Space direction="vertical" className="w-full" size={12}>
								<Alert type="success" showIcon message="Redis 已更新，后续百度网盘文件操作会读取新的 access token。" />
								<Descriptions
									size="small"
									column={1}
									items={[
										{ key: "expires_in", label: "有效期秒数", children: token.expires_in ?? "-" },
										{ key: "access_token", label: "Access Token", children: maskToken(token.access_token) },
										{ key: "refresh_token", label: "Refresh Token", children: maskToken(token.refresh_token) },
									]}
								/>
							</Space>
						)
						: <Text type="secondary">尚未保存新的 token。</Text>}
				</section>
			</div>
		</BasicContent>
	);
}
