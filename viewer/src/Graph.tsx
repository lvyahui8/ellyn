import { Graph } from '@antv/g6';
import { Button,Input,Space} from 'antd';
const { Search } = Input;
import { Col, Row } from 'antd';

import {useSearchParams} from 'react-router-dom'

import { ExtensionCategory, register } from '@antv/g6';
import { GNode, Group, Image, Rect, Text } from '@antv/g6-extension-react';

import axios from "axios";
import {useEffect, useState} from "react";
import {Simulate} from "react-dom/test-utils";
import load = Simulate.load;

register(ExtensionCategory.NODE, 'g', GNode);

const Node = ({ data, size }) => {
    const [width, height] = size;

    const { name, type, file, begin, end,covered_rate, covered_blocks,has_error,cost } = data;
    const color = has_error ? '#30BF78' : '#F4664A';
    const radius = 4;
    const methodLine = end.line - begin.line + 1
    const titleMap = {
        covered_rate: 'Covered',
        methodLine : 'Line',
        cost: 'Time',
    };


    const format = (category, value) => {
        if (category === 'covered_rate') return `${value.toFixed(2)}%`;
        if (category === 'cost') return `${value}s`;
        return value.toString();
    };

    const highlight = (category, value) => {
        if (category === 'success') {
            if (value >= 90) return 'green';
            if (value < 60) return 'red';
            return 'gray';
        }
        if (category === 'time') {
            if (value <= 10) return 'green';
            if (value >= 30) return 'red';
            return 'gray';
        }
        if (value >= 20) return 'red';
        if (value >= 5) return 'orange';
        return 'gray';
    };

    return (
        <Group>
            <Rect width={width} height={height} stroke={color} fill={'white'} radius={radius}>
                <Rect width={width} height={20} fill={color} radius={[radius, radius, 0, 0]}>
                    <Image
                        src={
                            type === 'module'
                                ? 'https://gw.alipayobjects.com/mdn/rms_8fd2eb/afts/img/A*0HC-SawWYUoAAAAAAAAAAABkARQnAQ'
                                : 'https://gw.alipayobjects.com/mdn/rms_8fd2eb/afts/img/A*sxK0RJ1UhNkAAAAAAAAAAABkARQnAQ'
                        }
                        x={2}
                        y={2}
                        width={16}
                        height={16}
                    />
                    <Text text={name} textBaseline="top" fill="#fff" fontSize={14} dx={20} dy={2} />
                </Rect>
                <Group transform="translate(5,40)">
                    {Object.entries({ covered_rate, cost, methodLine }).map(([key, value], index) => (
                        <Group key={index} transform={`translate(${(index * width) / 3}, 0)`}>
                            <Text text={titleMap[key]} fontSize={12} fill="gray" />
                            <Text text={format(key, value)} fontSize={12} dy={16} fill={highlight(key, value)} />
                        </Group>
                    ))}
                </Group>
            </Rect>
        </Group>
    );
};


function TrafficGraph() {
    const [searchParams] = useSearchParams();
    const [loading,setLoading] = useState(false)
    const [id ,setId] = useState(searchParams.get('id'))

    // 初始化图表实例
    const [graph] = useState(new Graph({
        container: 'container',
        data: {},
        node: {
            type: 'g',
            style: {
                size: [180, 60],
                component: (data) => <Node data={data} size={[180, 60]} />,
            },
        },
        behaviors: ['drag-element', 'zoom-canvas', 'drag-canvas'],
        plugins: [
            {
                type: 'toolbar',
                position: 'top-left',
                onClick: (item) => {

                },
                getItems: () => {
                    // G6 内置了 9 个 icon，分别是 zoom-in、zoom-out、redo、undo、edit、delete、auto-fit、export、reset
                    return [
                        { id: 'zoom-in', value: 'zoom-in' },
                        { id: 'zoom-out', value: 'zoom-out' },
                        { id: 'auto-fit', value: 'auto-fit' },
                        { id: 'export', value: 'export' },
                        { id: 'reset', value: 'reset' },
                    ];
                },
            },
        ],
    }));

    function setTrafficId(event) {
        setId(event.target.value)
    }

    console.log("draw graph")

    function loadGraph() {
        if (! id) {
            return
        }
        // 获取input中的id值
        console.log("query id :" + id)
        // 调用后端获取到图数据
        setLoading(true)
        axios.get('http://localhost:19898/traffic/detail?id=' +id)
            .then(resp => {
                // 设置到graph中
                console.log("loaded")
                console.log(resp.data)
                graph.setData(resp.data)
                setLoading(false)
                graph.render()
            })
            .catch(err => {
                console.log(err)
                setLoading(false)
            })
    }

    useEffect(() => {
        console.log('graph effect')
        if (id) {
            loadGraph()
        } else {
            graph.render()
        }
    },[])

    return (
        <>
           <Row>
               <Col span={24}>
                   <Space.Compact>
                       <Input defaultValue={searchParams.get('id')} onChange={setTrafficId}  placeholder={"输入流量id"} disabled={loading} />
                       <Button type="primary" onClick={loadGraph} disabled={loading}>查询</Button>
                   </Space.Compact>
               </Col>
           </Row>
            <br/>
            <Row>
                <Col span={24}>
                    <div id="container" style={{width: '100%',height: '700px', background:'#cccccc'}}></div>
                </Col>
            </Row>
        </>
    )
}

export default TrafficGraph