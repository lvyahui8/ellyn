import {Drawer} from "antd";
import {useContext, useEffect, useState} from "react";

import {graphCtx} from './Graph.tsx'
import axios from "axios";

import CodeMirror, {EditorView} from '@uiw/react-codemirror';
import { StreamLanguage } from '@codemirror/language';
import { go } from '@codemirror/legacy-modes/mode/go';
import {whiteLight} from '@uiw/codemirror-theme-white'

import { classname } from '@uiw/codemirror-extensions-classname';

const themeConf = EditorView.baseTheme({
    // '&dark .target-line': { backgroundColor: 'yellow' },
    '&light .covered-line': { backgroundColor: 'lightgreen' },
    '&light .uncovered-line': { backgroundColor: '#ff4d4f' },
});


const NodeDetail = () => {
    const {detailView,closeDetail,id,nodeId} = useContext(graphCtx)
    const [code,setCode] = useState("")
    const [data,setData] = useState()

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

    // 抽屉显示方法出入参数、代码行，覆盖明细，耗时，异常等信息
    return (
        <Drawer title="name" onClose={onClose}  open={detailView} getContainer={false} size={"large"}>
            <CodeMirror value={code} height="300px"
                        extensions={[StreamLanguage.define(go),classnameExt]}
                        theme={[whiteLight,themeConf]} editable={false}/>
        </Drawer>
    )
}

export default NodeDetail