import {
	BellOutlined,
	GlobalOutlined,
	LockOutlined,
	SafetyCertificateOutlined,
	UserOutlined,
} from "@ant-design/icons";

import { Card, List, Select, Space, Switch, Typography } from "antd";
import { BasicContent } from "#src/components/basic-content";

const { Paragraph, Text, Title } = Typography;

export default function Settings() {
	return (
		<BasicContent className="min-h-full">
			<Space direction="vertical" size={16} className="w-full max-w-5xl">
				<Card>
					<Space direction="vertical" size={4}>
						<Title level={3} className="!mb-0">设置</Title>
						<Paragraph type="secondary" className="!mb-0">
							管理账号偏好、安全策略和系统通知方式。
						</Paragraph>
					</Space>
				</Card>

				<Card title={(
					<Space>
						<UserOutlined />
						账号偏好
					</Space>
				)}
				>
					<List
						itemLayout="horizontal"
						dataSource={[
							{
								icon: <GlobalOutlined />,
								title: "界面语言",
								description: "选择后台默认显示语言。",
								action: <Select className="w-36" value="zh-CN" options={[{ label: "简体中文", value: "zh-CN" }, { label: "English", value: "en-US" }]} />,
							},
							{
								icon: <BellOutlined />,
								title: "系统提醒",
								description: "接收关键操作和任务状态提醒。",
								action: <Switch defaultChecked />,
							},
						]}
						renderItem={item => (
							<List.Item actions={[item.action]}>
								<List.Item.Meta
									avatar={item.icon}
									title={<Text strong>{item.title}</Text>}
									description={item.description}
								/>
							</List.Item>
						)}
					/>
				</Card>

				<Card title={(
					<Space>
						<SafetyCertificateOutlined />
						安全
					</Space>
				)}
				>
					<List
						itemLayout="horizontal"
						dataSource={[
							{
								icon: <LockOutlined />,
								title: "登录保护",
								description: "检测异常登录并提醒当前账号。",
								action: <Switch defaultChecked />,
							},
							{
								icon: <SafetyCertificateOutlined />,
								title: "敏感操作确认",
								description: "执行删除、授权等操作时进行二次确认。",
								action: <Switch defaultChecked />,
							},
						]}
						renderItem={item => (
							<List.Item actions={[item.action]}>
								<List.Item.Meta
									avatar={item.icon}
									title={<Text strong>{item.title}</Text>}
									description={item.description}
								/>
							</List.Item>
						)}
					/>
				</Card>
			</Space>
		</BasicContent>
	);
}
