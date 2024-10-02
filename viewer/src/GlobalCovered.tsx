
import {Col, Divider, Row, Statistic, Tree} from "antd"
import CodeMirror from "@uiw/react-codemirror";
import {StreamLanguage} from "@codemirror/language";
import {go} from "@codemirror/legacy-modes/mode/go";
import {whiteLight} from "@uiw/codemirror-theme-white";
import {ProDescriptions} from "@ant-design/pro-components";
import {useEffect, useState} from "react";
import axios from "axios";
import {themeConf} from "./NodeDetail.tsx";
import {classname} from "@uiw/codemirror-extensions-classname";

const { DirectoryTree } = Tree;


const GlobalCovered =  () => {
    const [data, setData] = useState([])

    const [fileInfo, setFileInfo] = useState({
        code : "",
        lineMap : {}
    })

    const [target,setTarget] = useState({
        totalLineNum : 0,
        targetLineNum  : 0,
        coveredLineNum : 0,
        coveredRate : 0.0
    })

    const loadCode = function (id) {
        axios.get('http://localhost:19898/source/file?id=' + id)
            .then(resp => {
                setFileInfo(resp.data)
            })
            .catch(err => {
                console.log(err)
            })
    }

    const onSelect = (keys, info) => {
        console.log('Trigger Select', keys, info);
        if (! info.node.isLeaf) {
            return
        }
        loadCode(keys[0])
    };
    const onExpand = (keys, info) => {
        console.log('Trigger Expand', keys, info);
    };

    const classnameExt = classname({
        add: (lineNumber) => {
            const flag = fileInfo.lineMap[lineNumber]
            if (flag) {
                if (flag == 1) {
                    return "uncovered-line"
                } else {
                    return "covered-line"
                }
            }
            return ""
        }
    })

    useEffect(() => {
        axios.get('http://localhost:19898/source/tree')
            .then(resp => {
                setData(resp.data)
                for (let i = 0; i < resp.data.length; i++) {
                    let n = resp.data[i]
                    if (n.isLeaf) {
                        loadCode(n.key)
                        break
                    }
                }
            })
            .catch(err => {
                console.log(err)
            })

        axios.get('http://localhost:19898/target/info')
            .then(resp => {
                setTarget(resp.data)
            })
            .catch(err => {
                console.log(err)
            })
    },[])

    return (
        <>
            <Row>
                <Col span={6}>
                    <Statistic title="总代码行" value={target.totalLineNum} />
                </Col>
                <Col span={6}>
                    <Statistic title="插桩代码行" value={target.targetLineNum} />
                </Col>
                <Col span={6}>
                    <Statistic title="已覆盖行" value={target.coveredLineNum} />
                </Col>
                <Col span={6}>
                    <Statistic title="覆盖率" value={target.coveredRate + "%"} />
                </Col>
            </Row>
            <Divider />
            <Row>
                <Col span={4}>
                    <DirectoryTree
                        multiple
                        defaultExpandAll
                        onSelect={onSelect}
                        onExpand={onExpand}
                        treeData={data}
                    />
                </Col>
                <Col span={20}>
                    <CodeMirror value={fileInfo.code} height="600px"
                                extensions={[StreamLanguage.define(go),classnameExt]}
                                theme={[whiteLight,themeConf]}
                                basicSetup={{
                                    highlightActiveLine : false
                                }}
                                editable={false}/>;
                </Col>
            </Row>
        </>
    );
}

export default GlobalCovered
