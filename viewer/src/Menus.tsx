import {useEffect, useState} from 'react';
import { AppstoreOutlined, MailOutlined, SettingOutlined } from '@ant-design/icons';
import type { MenuProps } from 'antd';
import { Menu } from 'antd';
import {useNavigate} from 'react-router-dom'


const Menus =  ({menuItems}) => {
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
        <Menu onClick={onClick} selectedKeys={[current]} mode="horizontal" items={menuItems} />
    )
};

export default Menus;