import Menus from "./Menus.tsx";

import {   Layout, theme  } from 'antd';
import {Route, Routes} from "react-router-dom";
import TrafficList from "./TrafficList.tsx";
import TrafficGraph from "./Graph.tsx";
import Meta from "./Meta.tsx";
import {ClusterOutlined, BarsOutlined,ProjectOutlined} from "@ant-design/icons";
const { Header, Content, Footer } = Layout;

const menuItems = [
    {
        label: '流量列表',
        key: '/traffic/list',
        icon: <BarsOutlined />,
        element : <TrafficList/>,
    },
    {
        label: '流量查询',
        key: '/traffic/query',
        icon:<ClusterOutlined />,
        element: <TrafficGraph/>,
    },
    {
        label: '元数据管理',
        key: '/meta',
        icon: <ProjectOutlined />,
        element: <Meta/>,
    },
]

const SiteLayout = () => {
    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();
    return (
        <Layout>
            <Header
                style={{
                    display: 'flex',
                    alignItems: 'center',
                }}
            >
                <div className="demo-logo" />
                <Menus menuItems={menuItems}/>
            </Header>
            <Content
                style={{
                    padding: '0 48px',
                }}
            >
                <div
                    style={{
                        background: colorBgContainer,
                        minHeight: 280,
                        padding: 24,
                        borderRadius: borderRadiusLG,
                    }}
                >
                    <Routes>
                        {
                            menuItems.map((item ) => (
                                <Route path={item.key} element={item.element} key={item.key}/>
                            ))
                        }
                    </Routes>
                </div>
            </Content>
            <Footer
                style={{
                    textAlign: 'center',
                }}
            >
                Ellyn ©{new Date().getFullYear()} Created by lvyahui8(Feego)
            </Footer>
        </Layout>
    );
};

export default SiteLayout;


