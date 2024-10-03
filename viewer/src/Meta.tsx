import {Button, Col, Row, Table} from 'antd';
import {useEffect, useState} from "react";
import axios from "axios";

function Meta() {
    const [data, setData] = useState(null)
    const [loading,setLoading] = useState(true)
    const [error, setError] = useState(null)

    useEffect(() => {
        axios.get('http://localhost:19898/meta/methods')
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

    const columns = [
        {
            title: '函数id',
            dataIndex: ['method','Id'],
        },
        {
            title: '函数名',
            dataIndex: 'method.FullName'.split('.'),
        },
        {
            title: '文件',
            dataIndex: 'file',
        },
        {
            title: '包',
            dataIndex: 'package',
        },
        {
            title: 'Block数',
            dataIndex: 'method.BlockCnt'.split('.'),
        },
        {
            title: '行数',
            render : function (text,record,index) {
                console.log(record)
                return record.method.End.line - record.method.Begin.line + 1
            }
        },
        {
            title: '参数列表',
            dataIndex: 'method.ArgsList'.split('.'),
        },
        {
            title: '返回值列表',
            dataIndex: 'method.ReturnList'.split('.'),
        },
        {
            title : '操作',
            render : function(text, record, index) {
                return <Button type={"primary"}   >配置mock</Button>
            }
        }
    ];
    return <Table dataSource={data} columns={columns} rowKey={"Id"} />
}

export default Meta