import type { ColProps } from "antd";
import {
	MessageOutlined,
	MoneyCollectOutlined,
	ShoppingCartOutlined,
	UserOutlined,
} from "@ant-design/icons";
import { Button, Card, Col, Row, Statistic, Typography } from "antd";
import CountUp from "react-countup";
import { useTranslation } from "react-i18next";

const { Text } = Typography;

const wrapperCol: ColProps = {
	xs: 24,
	sm: 24,
	md: 12,
	lg: 12,
	xl: 12,
	xxl: 6,
};

export default function CardList() {
	const { t } = useTranslation();

	const CARD_LIST = [
		{
			title: t("home.newVisits"),
			data: 102_400,
			icon: <UserOutlined />,
		},
		{
			title: t("home.messages"),
			data: 81212,
			icon: <MessageOutlined />,
		},
		{
			title: t("home.purchases"),
			data: 9280,
			icon: <MoneyCollectOutlined />,
		},
		{
			title: t("home.shoppings"),
			data: 13600,
			icon: <ShoppingCartOutlined />,
		},

	];

	return (
		<Row justify="space-between" gutter={[20, 20]}>
			{
				CARD_LIST.map((cardItem) => {
					return (
						<Col key={cardItem.title} {...wrapperCol}>
							<Card>
								<div className="flex items-center justify-between gap-4">
									<Statistic
										title={<Text type="secondary">{cardItem.title}</Text>}
										value={cardItem.data}
										formatter={value => <CountUp end={Number(value)} separator="," />}
									/>
									<Button
										className="text-2xl"
										icon={cardItem.icon}
										type="text"
									/>
								</div>
							</Card>
						</Col>
					);
				})
			}
		</Row>
	);
}
