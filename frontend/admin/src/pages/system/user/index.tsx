import type { ColumnsType } from "antd/es/table";
import type { Dayjs } from "dayjs";
import type { CreateUserPayload, DocumentRecord, UserRecord } from "#src/api/user";

import {
	FileTextOutlined,
	PlusOutlined,
	ReloadOutlined,
	SearchOutlined,
	UserOutlined,
} from "@ant-design/icons";
import {
	Button,
	DatePicker,
	Empty,
	Form,
	Input,
	message,
	Modal,
	Space,
	Table,
	Tag,
	Typography,
} from "antd";
import dayjs from "dayjs";
import { useEffect, useMemo, useState } from "react";
import { createUser, fetchUserDocuments, fetchUsers } from "#src/api/user";
import { BasicContent } from "#src/components/basic-content";

const { Text, Title } = Typography;

interface UserSearchForm {
	username?: string
	email?: string
	phone?: string
}

interface DocumentSearchForm {
	display_name?: string
	file_name?: string
}

type CreateUserForm = Omit<CreateUserPayload, "birthday"> & {
	birthday?: Dayjs
};

function formatTime(value?: string) {
	if (!value) {
		return "-";
	}
	return dayjs(value).format("YYYY-MM-DD HH:mm");
}

export default function User() {
	const [searchForm] = Form.useForm<UserSearchForm>();
	const [documentSearchForm] = Form.useForm<DocumentSearchForm>();
	const [createForm] = Form.useForm<CreateUserForm>();
	const [users, setUsers] = useState<UserRecord[]>([]);
	const [documents, setDocuments] = useState<DocumentRecord[]>([]);
	const [selectedUser, setSelectedUser] = useState<UserRecord>();
	const [usersLoading, setUsersLoading] = useState(false);
	const [documentsLoading, setDocumentsLoading] = useState(false);
	const [createOpen, setCreateOpen] = useState(false);
	const [creating, setCreating] = useState(false);

	const selectedUserId = selectedUser?.id;

	const loadUsers = async (values?: UserSearchForm) => {
		setUsersLoading(true);
		try {
			const nextUsers = await fetchUsers({
				...(values ?? {}),
				page: 1,
				page_size: 100,
			});
			setUsers(nextUsers);
			if (!selectedUserId && nextUsers.length > 0) {
				setSelectedUser(nextUsers[0]);
			}
			else if (selectedUserId && !nextUsers.some(user => user.id === selectedUserId)) {
				setSelectedUser(nextUsers[0]);
			}
		}
		finally {
			setUsersLoading(false);
		}
	};

	const loadDocuments = async (
		userId = selectedUserId,
		values?: DocumentSearchForm,
	) => {
		if (!userId) {
			setDocuments([]);
			return;
		}

		setDocumentsLoading(true);
		try {
			const nextDocuments = await fetchUserDocuments({
				user_id: userId,
				...(values ?? {}),
				page: 1,
				page_size: 100,
			});
			setDocuments(nextDocuments);
		}
		finally {
			setDocumentsLoading(false);
		}
	};

	useEffect(() => {
		loadUsers();
	}, []);

	useEffect(() => {
		loadDocuments(selectedUserId);
	}, [selectedUserId]);

	const userColumns = useMemo<ColumnsType<UserRecord>>(() => [
		{
			title: "用户",
			dataIndex: "username",
			key: "username",
			render: (value: string, record: UserRecord) => (
				<Space orientation="vertical" size={0}>
					<Text strong>{value}</Text>
					<Text type="secondary" className="text-xs">{record.email}</Text>
				</Space>
			),
		},
		{
			title: "手机号",
			dataIndex: "phone",
			key: "phone",
			width: 140,
			render: (value: string) => value || "-",
		},
		{
			title: "生日",
			dataIndex: "birthday",
			key: "birthday",
			width: 120,
			render: (value: string) => value ? dayjs(value).format("YYYY-MM-DD") : "-",
		},
		{
			title: "创建时间",
			dataIndex: "created_at",
			key: "created_at",
			width: 170,
			render: formatTime,
		},
	], []);

	const documentColumns = useMemo<ColumnsType<DocumentRecord>>(() => [
		{
			title: "Document",
			dataIndex: "display_name",
			key: "display_name",
			render: (value: string, record: DocumentRecord) => (
				<Space orientation="vertical" size={0}>
					<Text strong>{value}</Text>
					<Text type="secondary" className="text-xs">{record.file_name || "无文件名"}</Text>
				</Space>
			),
		},
		{
			title: "ID",
			dataIndex: "id",
			key: "id",
			width: 90,
			render: (value: number) => <Tag>{value}</Tag>,
		},
		{
			title: "创建时间",
			dataIndex: "created_at",
			key: "created_at",
			width: 170,
			render: formatTime,
		},
		{
			title: "更新时间",
			dataIndex: "updated_at",
			key: "updated_at",
			width: 170,
			render: formatTime,
		},
	], []);

	const handleCreateUser = async () => {
		const values = await createForm.validateFields();
		setCreating(true);
		try {
			await createUser({
				...values,
				birthday: values.birthday?.format("YYYY-MM-DD"),
			});
			message.success("用户已创建");
			setCreateOpen(false);
			createForm.resetFields();
			await loadUsers();
		}
		finally {
			setCreating(false);
		}
	};

	return (
		<BasicContent className="min-h-full">
			<div className="flex flex-col gap-4">
				<div className="flex flex-wrap items-center justify-between gap-3">
					<div>
						<Title level={3} className="mb-1!">用户管理</Title>
						<Text type="secondary">创建用户，选择用户后查看对应 documents。</Text>
					</div>
					<Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateOpen(true)}>
						创建用户
					</Button>
				</div>

				<div className="grid grid-cols-1 gap-4 xl:grid-cols-[minmax(560px,1fr)_minmax(460px,0.85fr)]">
					<section className="rounded-md border border-colorBorder bg-colorBgContainer p-4">
						<div className="mb-3 flex items-center gap-2">
							<UserOutlined />
							<Text strong>用户列表</Text>
						</div>
						<Form
							form={searchForm}
							layout="inline"
							className="mb-4 gap-y-2"
							onFinish={loadUsers}
						>
							<Form.Item name="username">
								<Input allowClear placeholder="用户名" />
							</Form.Item>
							<Form.Item name="email">
								<Input allowClear placeholder="邮箱" />
							</Form.Item>
							<Form.Item name="phone">
								<Input allowClear placeholder="手机号" />
							</Form.Item>
							<Form.Item>
								<Space>
									<Button htmlType="submit" type="primary" icon={<SearchOutlined />}>查询</Button>
									<Button
										icon={<ReloadOutlined />}
										onClick={() => {
											searchForm.resetFields();
											loadUsers({});
										}}
									>
										重置
									</Button>
								</Space>
							</Form.Item>
						</Form>
						<Table<UserRecord>
							rowKey="id"
							size="middle"
							loading={usersLoading}
							columns={userColumns}
							dataSource={users}
							pagination={false}
							rowClassName={record => record.id === selectedUserId ? "ant-table-row-selected" : ""}
							onRow={record => ({
								onClick: () => setSelectedUser(record),
							})}
						/>
					</section>

					<section className="rounded-md border border-colorBorder bg-colorBgContainer p-4">
						<div className="mb-3 flex items-center justify-between gap-3">
							<Space>
								<FileTextOutlined />
								<Text strong>用户 Documents</Text>
							</Space>
							{selectedUser ? <Tag color="blue">{selectedUser.username}</Tag> : null}
						</div>

						{selectedUser
							? (
								<>
									<Form
										form={documentSearchForm}
										layout="inline"
										className="mb-4 gap-y-2"
										onFinish={values => loadDocuments(selectedUser.id, values)}
									>
										<Form.Item name="display_name">
											<Input allowClear placeholder="Document 名称" />
										</Form.Item>
										<Form.Item name="file_name">
											<Input allowClear placeholder="文件名" />
										</Form.Item>
										<Form.Item>
											<Space>
												<Button htmlType="submit" type="primary" icon={<SearchOutlined />}>查询</Button>
												<Button
													icon={<ReloadOutlined />}
													onClick={() => {
														documentSearchForm.resetFields();
														loadDocuments(selectedUser.id, {});
													}}
												>
													重置
												</Button>
											</Space>
										</Form.Item>
									</Form>
									<Table<DocumentRecord>
										rowKey="id"
										size="middle"
										loading={documentsLoading}
										columns={documentColumns}
										dataSource={documents}
										pagination={false}
									/>
								</>
							)
							: <Empty description="先选择一个用户" />}
					</section>
				</div>
			</div>

			<Modal
				title="创建用户"
				open={createOpen}
				onCancel={() => setCreateOpen(false)}
				onOk={handleCreateUser}
				confirmLoading={creating}
				okText="创建"
				cancelText="取消"
				destroyOnHidden
			>
				<Form form={createForm} layout="vertical" className="pt-2">
					<Form.Item label="用户名" name="username" rules={[{ required: true, message: "请输入用户名" }]}>
						<Input placeholder="例如 xiaofeng" />
					</Form.Item>
					<Form.Item
						label="邮箱"
						name="email"
						rules={[
							{ required: true, message: "请输入邮箱" },
							{ type: "email", message: "邮箱格式不正确" },
						]}
					>
						<Input placeholder="name@example.com" />
					</Form.Item>
					<Form.Item label="密码" name="password" rules={[{ required: true, min: 6, message: "请输入至少 6 位密码" }]}>
						<Input.Password placeholder="至少 6 位" />
					</Form.Item>
					<Form.Item label="手机号" name="phone">
						<Input placeholder="可选" />
					</Form.Item>
					<Form.Item label="生日" name="birthday">
						<DatePicker className="w-full" />
					</Form.Item>
					<Form.Item label="User Prompt" name="user_prompt">
						<Input.TextArea rows={4} placeholder="可选，用户默认提示词" />
					</Form.Item>
				</Form>
			</Modal>
		</BasicContent>
	);
}
