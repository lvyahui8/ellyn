import {Card,Col,Row,Typography,Button,Space} from "antd";

import axios from "axios";
import {useLocation, useNavigate} from 'react-router-dom'
import {useState} from "react";


const { Text, Link } = Typography;

interface CaseRes {
    gid : string
    resp : string
}

interface Case {
    title: string,
    url: string,
    method: string,
    body? : string,
    res? : CaseRes,
}

const TestCase = () => {
    const navigate = useNavigate()


    const [cases,setCases] = useState<Case[]>([
        {
            "title" : "测试覆盖率采集",
            "url" : "http://localhost:8080/user/foo",
            "method" : "GET",
            "body" : "",
        },
        {
            "title" : "测试异步链路采集&合并",
            "url" : "http://localhost:8080/trade",
            "method" : "GET",
            "body" : "",
        },
        {
            "title" : "测试参数采集",
            "url" : "http://localhost:8080/profile/1",
            "method" : "GET",
            "body" : "",
        },
        {
            "title" : "测试参数采集",
            "url" : "http://localhost:8080/profile/1",
            "method" : "GET",
            "body" : "",
        },
        {
            "title" : "测试参数采集",
            "url" : "http://localhost:8080/profile/1",
            "method" : "GET",
            "body" : "",
        },
    ])

    const runTestCase = function (item,idx) {
        console.log("run testcase :" + idx)
        axios.get(item.url).then(resp => {
            item.res = {
                resp: JSON.stringify(resp.data)
            }
            setCases(cases.map(item=>item))
        })
    }

    return (
        <>
            <Row gutter={[16,16]}>
            {
                cases.map((item,idx) => {
                    return <Col span={8}>
                        <Card
                            title={item.title}
                        >
                            <Space direction={"vertical"}>
                                <Text code>GET {item.url}</Text>
                                <Button onClick={(e) => {runTestCase(item,idx)}} type={"primary"} >测试请求</Button>
                                {
                                    item.res && (
                                        <>
                                            Graph:<Link href="https://ant.design" target="_blank"> </Link>
                                            Response:
                                            <Text code>
                                                {item.res.resp}
                                            </Text>
                                        </>
                                    )
                                }
                            </Space>
                        </Card>
                    </Col>
                })
            }
            </Row>
        </>
    )
}

export default TestCase