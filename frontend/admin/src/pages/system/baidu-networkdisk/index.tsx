import type { BaiduTokenResponse } from "#src/api/baidu-networkdisk";

import {
	CloudOutlined,
	ExportOutlined,
	KeyOutlined,
	SaveOutlined,
} from "@ant-design/icons";
import {
	Alert,
	Button,
	Card,
	Descriptions,
	Flex,
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

const baiduAuthorizeUrl = "https://openapi.baidu.com/oauth/2.0/authorize?response_type=code&client_id=AhdQdE8PIYnUYSFKOMm7LBwIbpqaZpCE&redirect_uri=oob&scope=basic,netdisk";

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

	const openAuthorizePage = () => {
		window.open(baiduAuthorizeUrl, "_blank", "noopener,noreferrer");
	};

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
			<Space direction="vertical" size={16} className="w-full max-w-4xl">
				<Card>
					<Title level={3} className="!mb-1">百度网盘授权</Title>
					<Text type="secondary">先打开百度授权页获取 code，再在这里换取 token 并写入 Redis。</Text>
				</Card>

				<Card>
					<Flex className="mb-4" wrap gap={12} align="center" justify="space-between">
						<Space>
							<CloudOutlined />
							<Text strong>获取授权 code</Text>
						</Space>
						<Button icon={<ExportOutlined />} onClick={openAuthorizePage}>
							打开授权页
						</Button>
					</Flex>
					<Alert
						type="info"
						showIcon
						className="mb-4"
						message="授权页会在新窗口打开，复制页面返回的 code 后粘贴到下方输入框。"
					/>
					<div className="mb-4 flex items-center gap-2">
						<KeyOutlined />
						<Text strong>填写 code 并获取 token</Text>
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
								获取 token
							</Button>
						</Form.Item>
					</Form>
				</Card>

				<Card>
					<Space className="mb-4">
						<CloudOutlined />
						<Text strong>保存结果</Text>
					</Space>
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
				</Card>
			</Space>
		</BasicContent>
	);
}
