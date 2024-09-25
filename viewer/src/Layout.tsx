import Menus from "./Menus.tsx";

import {   Layout, theme  } from 'antd';
import {Route, Routes} from "react-router-dom";
import TrafficList from "./TrafficList.tsx";
import TrafficGraph from "./Graph.tsx";
import Meta from "./Meta.tsx";
const { Header, Content, Footer } = Layout;


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
                <Menus/>
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
                        <Route path = '/traffic/list' element = {<TrafficList/>} />
                        <Route path = '/traffic/query' element = {<TrafficGraph/>} />
                        <Route path = '/meta' element = {<Meta/>} />
                    </Routes>
                </div>
            </Content>
            <Footer
                style={{
                    textAlign: 'center',
                }}
            >
                Ant Design ©{new Date().getFullYear()} Created by Ant UED
            </Footer>
        </Layout>
    );
};

export default SiteLayout;


