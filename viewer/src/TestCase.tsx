import {Card,Col,Row} from "antd";


const TestCase = () => {
    return (
        <>
            <Row gutter={16}>
                <Col span={8}>
                    <Card
                        title="测试覆盖率采集"
                    >
                    </Card>
                </Col>
                <Col span={8}>
                    <Card
                        title="测试异步链路采集&合并"
                    >
                    </Card>
                </Col>
                <Col span={8}>
                    <Card
                        title="测试参数采集"
                    >
                    </Card>
                </Col>
            </Row>
        </>
    )
}

export default TestCase