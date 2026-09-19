import type { UserRecord } from "#src/api/user";

import { addUserPoints, deductUserPoints, fetchUserPoints } from "#src/api/points";
import {
	FallOutlined,
	RiseOutlined,
} from "@ant-design/icons";
import {
	Button,
	Empty,
	Form,
	Input,
	InputNumber,
	Modal,
	Space,
	Spin,
	Tag,
	Typography,
	message,
} from "antd";
import { useEffect, useState } from "react";
import { createUseStyles } from "react-jss";

const { Text } = Typography;

type AdjustAction = "add" | "deduct";

interface AdjustForm {
	amount: number
	remark: string
}

interface UserPointsPanelProps {
	user?: UserRecord
}

const useStyles = createUseStyles(({ token, isDark }) => ({
	panel: {
		position: "relative",
		overflow: "hidden",
		border: `1px solid ${token.colorBorderSecondary}`,
		borderRadius: 18,
		background: isDark
			? `linear-gradient(180deg, ${token.colorBgElevated} 0%, ${token.colorBgContainer} 100%)`
			: "linear-gradient(180deg, #f6f1e4 0%, #fffaf3 52%, #ffffff 100%)",
		boxShadow: isDark ? token.boxShadowTertiary : "0 18px 40px -28px rgba(92, 58, 18, 0.45)",
		"&::before": {
			content: '""',
			position: "absolute",
			inset: "0 0 auto",
			height: 3,
			background: "linear-gradient(90deg, #b45309 0%, #d97706 42%, #f4d38a 100%)",
		},
	},
	inner: {
		padding: "18px 18px 16px",
	},
	meta: {
		display: "flex",
		alignItems: "center",
		justifyContent: "space-between",
		gap: 12,
		marginBottom: 10,
	},
	kicker: {
		margin: 0,
		color: token.colorTextSecondary,
		fontSize: 11,
		fontWeight: 700,
		letterSpacing: "0.18em",
		textTransform: "uppercase",
	},
	balance: {
		margin: "4px 0 0",
		color: isDark ? token.colorText : "#3b2a14",
		fontFamily: 'ui-monospace, "SFMono-Regular", "Cascadia Code", Menlo, monospace',
		fontSize: 42,
		fontWeight: 700,
		lineHeight: 1,
		fontVariantNumeric: "tabular-nums",
		fontFeatureSettings: '"tnum"',
	},
	actions: {
		display: "flex",
		gap: 8,
		marginTop: 16,
	},
	credit: {
		flex: 1,
	},
	debit: {
		flex: 1,
	},
	"@media (max-width: 768px)": {
		inner: {
			padding: "16px 14px 14px",
		},
		meta: {
			alignItems: "flex-start",
			flexDirection: "column",
		},
		balance: {
			fontSize: 34,
		},
		actions: {
			flexDirection: "column",
		},
	},
	"@supports (backdrop-filter: blur(8px))": {
		panel: {
			backdropFilter: "blur(8px)",
		},
	},
}));

export function UserPointsPanel({ user }: UserPointsPanelProps) {
	const classes = useStyles();
	const [form] = Form.useForm<AdjustForm>();
	const [balance, setBalance] = useState<number>();
	const [loading, setLoading] = useState(false);
	const [submitting, setSubmitting] = useState(false);
	const [action, setAction] = useState<AdjustAction>("add");
	const [open, setOpen] = useState(false);

	const userId = user?.id;
	const isDeduct = action === "deduct";

	const loadBalance = async (id = userId) => {
		if (!id) {
			setBalance(undefined);
			return;
		}
		setLoading(true);
		try {
			const result = await fetchUserPoints(id);
			setBalance(result.balance);
		}
		finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		loadBalance(userId);
	}, [userId]);

	const openAdjust = (nextAction: AdjustAction) => {
		setAction(nextAction);
		form.setFieldsValue({
			amount: 1,
			remark: nextAction === "add" ? "管理员增加积分" : "管理员扣减积分",
		});
		setOpen(true);
	};

	const handleSubmit = async () => {
		if (!userId) {
			return;
		}
		const values = await form.validateFields();
		setSubmitting(true);
		try {
			const result = isDeduct
				? await deductUserPoints(userId, values)
				: await addUserPoints(userId, values);
			setBalance(result.balance);
			setOpen(false);
			form.resetFields();
			message.success(isDeduct ? "积分已扣减" : "积分已增加");
		}
		finally {
			setSubmitting(false);
		}
	};

	return (
		<div className={classes.panel}>
			<div className={classes.inner}>
				{user
					? (
						<Spin spinning={loading}>
							<div className={classes.meta}>
								<div>
									<p className={classes.kicker}>Account Ledger</p>
									<Text type="secondary">{user.username} 的当前积分</Text>
								</div>
								<Tag color="gold">UID {user.id}</Tag>
							</div>
							<p className={classes.balance}>{balance ?? 0}</p>
							<div className={classes.actions}>
								<Button
									className={classes.credit}
									type="primary"
									icon={<RiseOutlined />}
									onClick={() => openAdjust("add")}
								>
									增加积分
								</Button>
								<Button
									className={classes.debit}
									danger
									icon={<FallOutlined />}
									disabled={(balance ?? 0) <= 0}
									onClick={() => openAdjust("deduct")}
								>
									减少积分
								</Button>
							</div>
						</Spin>
					)
					: <Empty description="先选择一个用户再调整积分" />}
			</div>

			<Modal
				title={isDeduct ? `扣减 ${user?.username ?? ""} 的积分` : `增加 ${user?.username ?? ""} 的积分`}
				open={open}
				onCancel={() => setOpen(false)}
				onOk={handleSubmit}
				confirmLoading={submitting}
				okText={isDeduct ? "确认扣减" : "确认增加"}
				okButtonProps={{ danger: isDeduct }}
				cancelText="取消"
				destroyOnHidden
			>
				<Form form={form} layout="vertical" className="pt-2">
					<Form.Item
						label="变动数量"
						name="amount"
						rules={[
							{ required: true, message: "请输入积分数量" },
							{
								type: "number",
								min: 1,
								message: "积分数量必须大于 0",
							},
						]}
					>
						<InputNumber
							className="w-full"
							min={1}
							max={isDeduct ? Math.max(balance ?? 1, 1) : undefined}
							precision={0}
							placeholder="正整数"
						/>
					</Form.Item>
					<Form.Item
						label="备注"
						name="remark"
						rules={[
							{ required: true, message: "请填写调整原因" },
							{ max: 255, message: "备注不能超过 255 个字符" },
						]}
					>
						<Input.TextArea rows={3} maxLength={255} showCount placeholder="记一笔调整原因，方便对账" />
					</Form.Item>
					<Space>
						<Text type="secondary">当前余额</Text>
						<Text strong>{balance ?? 0}</Text>
					</Space>
				</Form>
			</Modal>
		</div>
	);
}
