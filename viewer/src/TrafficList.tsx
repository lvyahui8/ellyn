
import {Col, MenuProps, Row, Table} from 'antd';
import { Button,Input,Space} from 'antd';

import {useEffect, useState} from "react";
import axios from "axios";
import {useLocation, useNavigate} from 'react-router-dom'


function TrafficList() {
    const [data, setData] = useState(null)
    const [loading,setLoading] = useState(true)
    const [error, setError] = useState(null)
    const navigate = useNavigate()

    useEffect(() => {
        axios.get('http://localhost:19898/traffic/list')
            .then(resp => {
                setData(resp.data)
                setLoading(false)
            })
            .catch(err => {
                setError(err.message)
                setLoading(false)
            })
    },[])

    if (loading) {
        return <div>Loading...</div>
    }
    if (error) {
        return <div>Error: {error}</div>
    }

    const onClick: MenuProps['onClick'] = (e, target) => {
        console.log("click to " + target)
        navigate('/traffic/query?id=' + target)
    };

    const columns = [
        {
            title: '流量id',
            dataIndex: 'id',
        },
        {
            title : '发生时间',
            dataIndex: 'time',
        },
        {
            title: '节点数量',
            render : function(text, record, index) {
                if (record.nodes == null) {
                    return 0
                }
                return record.nodes.length
            },
        },
        {
            title: '边数量',
            render : function(text, record, index) {
                if (record.edges == null) {
                    return 0
                }
                return record.edges.length
            },
        },
        {
            title : '操作',
            render : function(text, record, index) {
                return <Button type={"primary"}  data={"/traffic/detail/" + record.id} onClick={(e) => onClick(e,record.id)} >查看</Button>
            }
        }
    ];
    return <Table dataSource={data} columns={columns} rowKey={"id"} />
}

export default TrafficList