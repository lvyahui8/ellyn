import {useEffect, useState} from 'react';
import { AppstoreOutlined, MailOutlined, SettingOutlined } from '@ant-design/icons';
import type { MenuProps } from 'antd';
import { Menu } from 'antd';
import {useNavigate} from 'react-router-dom'


type MenuItem = Required<MenuProps>['items'][number];

const items: MenuItem[] = [
    {
        label: '流量列表',
        key: '/traffic/list',
        icon: <MailOutlined />,
    },
    {
        label: '流量查询',
        key: '/traffic/query',
        icon: <AppstoreOutlined />,
    },
    {
        label: '元数据管理',
        key: '/meta',
        icon: <AppstoreOutlined />,
    },
];

const Menus =  () => {
    const [current, setCurrent] = useState('/traffic/list');
    const navigate = useNavigate()

    useEffect(() => {
        console.log("to " + current)
        navigate(current)
    },[current])

    const onClick: MenuProps['onClick'] = (e) => {
        setCurrent(e.key);
    };

    return (
        <Menu onClick={onClick} selectedKeys={[current]} mode="horizontal" items={items} />
    )
};

export default Menus;