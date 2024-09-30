
import {Col, Divider, Row, Statistic, Tree} from "antd"
import CodeMirror from "@uiw/react-codemirror";
import {StreamLanguage} from "@codemirror/language";
import {go} from "@codemirror/legacy-modes/mode/go";
import {whiteLight} from "@uiw/codemirror-theme-white";
import {ProDescriptions} from "@ant-design/pro-components";

const { DirectoryTree } = Tree;

const treeData = [
    {
        title: 'parent 0',
        key: '0-0',
        children: [
            {
                title: 'leaf 0-0',
                key: '0-0-0',
                isLeaf: true,
            },
            {
                title: 'leaf 0-1',
                key: '0-0-1',
                isLeaf: true,
            },
        ],
    },
    {
        title: 'parent 1',
        key: '0-1',
        children: [
            {
                title: 'leaf 1-0',
                key: '0-1-0',
                isLeaf: true,
            },
            {
                title: 'leaf 1-1',
                key: '0-1-1',
                isLeaf: true,
            },
        ],
    },
];


const GlobalCovered =  () => {
    const onSelect = (keys, info) => {
        console.log('Trigger Select', keys, info);
    };
    const onExpand = (keys, info) => {
        console.log('Trigger Expand', keys, info);
    };
    return (
        <>
            <Row>
                <Col span={6}>
                    <Statistic title="总代码行" value={112893} />
                </Col>
                <Col span={6}>
                    <Statistic title="插桩代码行" value={112893} />
                </Col>
                <Col span={6}>
                    <Statistic title="已覆盖行" value={112893} />
                </Col>
                <Col span={6}>
                    <Statistic title="覆盖率" value={"23.6%"} />
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
                        treeData={treeData}
                    />
                </Col>
                <Col span={20}>
                    <CodeMirror value={"test"} height="600px" extensions={[StreamLanguage.define(go)]} theme={whiteLight} editable={false}/>;
                </Col>
            </Row>
        </>
    );
}

export default GlobalCovered
