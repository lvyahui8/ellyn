import {Drawer, Table} from "antd";
import {useContext, useEffect, useState} from "react";
import { ProDescriptions } from '@ant-design/pro-components';


import {graphCtx} from './Graph.tsx'
import axios from "axios";

import CodeMirror, {EditorView} from '@uiw/react-codemirror';
import { lineNumbers } from "@codemirror/view";
import { StreamLanguage } from '@codemirror/language';
import { go } from '@codemirror/legacy-modes/mode/go';
import {whiteLight} from '@uiw/codemirror-theme-white'

import { classname } from '@uiw/codemirror-extensions-classname';

import ReactJson from '@microlink/react-json-view'


export const themeConf = EditorView.baseTheme({
    // '&dark .target-line': { backgroundColor: 'yellow' },
    '&light .covered-line': { backgroundColor: 'lightgreen' },
    '&light .uncovered-line': { backgroundColor: '#ff8b8d' },
    '.cm-content' : {fontFamily : "consolas, Monaco, Lucida Console, monospace"},
});


const NodeDetail = () => {
    const {detailView,closeDetail,id,nodeId} = useContext(graphCtx)
    const [code,setCode] = useState("")
    const [data,setData] = useState<any>()
    const [lineNumberOffset, setLineNumberOffset] = useState(0)

    const onClose = () => {
        console.log("close trigger")
        closeDetail()
    };


    useEffect(() => {
        console.log("refresh Drawer")
        if (nodeId == "-1") {
            return
        }
        axios.get('http://localhost:19898/node/detail?graphId=' + id + "&nodeId=" + nodeId)
            .then(resp => {
                console.log(resp.data)
                setLineNumberOffset(resp.data.resNode.begin.line)
                setData(resp.data)
                setCode(resp.data.funcCode)
            })
            .catch(err => {
                console.log(err.message)
            })
    },[nodeId])

    const classnameExt = classname({
        add: (lineNumber) => {
            if (!data) {
                return
            }
            for (let i = 0; i < data.resNode.covered_blocks.length; i++) {
                const block = data.resNode.covered_blocks[i]
                const line = lineNumber + data.resNode.begin.line - 1
                if (line >= block.begin.line && line <= block.end.line) {
                    return 'covered-line';
                }
            }
            return 'uncovered-line';
        },
    });

    const lineNumberExt =  lineNumbers({
        formatNumber: (n, s) => {
            return (n + lineNumberOffset - 1).toString();
        }
    })
    const columns = [
        {
            title: '索引',
            dataIndex: 'idx',
        },
        {
            title: '变量名',
            dataIndex: 'name',
        },
        {
            title: '变量类型',
            dataIndex: 'type',
        },
        {
            title: '变量值',
            // dataIndex: 'val',
            render : function(text,record,index) {
                let res = JSON.parse(record.val)
                if (res == null) {
                    return "null"
                }
                if (typeof res === "object") {
                    return (
                        <ReactJson src={res}/>
                    )
                } else {
                    return res
                }
            }
        },
    ]
    // 抽屉显示方法出入参数、代码行，覆盖明细，耗时，异常等信息
    return (
        <Drawer title={data && data.resNode.name} onClose={onClose}  open={detailView} getContainer={false} size={"large"}>
            {
                data &&
                <>
                    <ProDescriptions
                        column={2}
                    >
                        <ProDescriptions.Item
                            label="函数名"
                            valueType="text">
                            {data.resNode.name}
                        </ProDescriptions.Item>
                        <ProDescriptions.Item
                            label="所在文件"
                            valueType="text">
                            {data.resNode.file}
                        </ProDescriptions.Item>
                        <ProDescriptions.Item
                            label="执行耗时"
                            valueType="text">
                            {data.resNode.cost}ms
                        </ProDescriptions.Item>
                        <ProDescriptions.Item
                            label="是否发生错误"
                            valueType="text">
                            {data.resNode.has_err ?
                                "是"
                                :
                                "否"}
                        </ProDescriptions.Item>

                        <ProDescriptions.Item
                            label="覆盖率"
                            valueType="text">
                            {data.resNode.covered_rate}%
                        </ProDescriptions.Item>

                    </ProDescriptions>
                    <br/>
                    <ProDescriptions
                        column={2}
                        layout={"vertical"}
                    >
                        <ProDescriptions.Item
                            label="参数列表"
                            span={2}
                            valueType="text">
                            <Table dataSource={data.resNode.args} size={"small"} pagination={false}
                                   columns={columns} rowKey={"idx"} locale={{emptyText:"无参数"}} />
                        </ProDescriptions.Item>
                        <ProDescriptions.Item
                            label="返回值列表"
                            span={2}
                            valueType="text">
                            <Table dataSource={data.resNode.returns} size={"small"} pagination={false}
                                   columns={columns} rowKey={"idx"} locale={{emptyText:"无参数"}} />
                        </ProDescriptions.Item>
                        <ProDescriptions.Item
                            label="覆盖明细"
                            span={2}
                            valueType="text">
                            <div style={{"width" : "100%"}}>
                                <CodeMirror value={code} height="300px"
                                            extensions={[ StreamLanguage.define(go), classnameExt,lineNumberExt]}
                                            theme={[whiteLight,themeConf]}
                                            basicSetup={{
                                                highlightActiveLine : false
                                            }}
                                            editable={false}/>
                            </div>
                        </ProDescriptions.Item>
                    </ProDescriptions>
                </>

            }
        </Drawer>
    )
}

export default NodeDetail