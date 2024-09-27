import {Drawer} from "antd";
import {useContext, useState} from "react";

import {graphCtx} from './Graph.tsx'

const NodeDetail = () => {
    const {detailView,closeDetail} = useContext(graphCtx)
    const onClose = () => {
        console.log("close trigger")
        closeDetail()
    };


    // 抽屉显示方法出入参数、代码行，覆盖明细，耗时，异常等信息
    return (
        <Drawer title="name" onClose={onClose}  open={detailView} getContainer={false} size={"large"}>

        </Drawer>
    )
}

export default NodeDetail