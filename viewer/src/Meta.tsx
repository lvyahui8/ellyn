import { Col, Row, Table } from 'antd';
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
            dataIndex: 'Id',
        },
        {
            title: '函数名',
            dataIndex: 'FullName',
        },
        {
            title: 'Block数',
            dataIndex: 'BlockCnt',
        },
    ];
    return <Table dataSource={data} columns={columns} rowKey={"Id"} />
}

export default Meta