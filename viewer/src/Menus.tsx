import {useEffect, useState} from 'react';
import { AppstoreOutlined, MailOutlined, SettingOutlined } from '@ant-design/icons';
import type { MenuProps } from 'antd';
import { Menu } from 'antd';
import {useLocation, useNavigate} from 'react-router-dom'
import {c} from "@codemirror/legacy-modes/mode/clike";


const Menus =  ({menuItems}) => {
    const location = useLocation()
    const navigate = useNavigate()
    const [current, setCurrent] = useState('')


    useEffect(() => {
        let targetPath = location.pathname
        if (! targetPath || targetPath == "/") {
            targetPath = String(menuItems[0].key)
            navigate(targetPath)
        }

    }, []);

    const onClick: MenuProps['onClick'] = (e) => {
        console.log("menu to " + e.key)
        navigate(e.key)
        setCurrent(e.key)
    };

    return (
        <Menu onClick={onClick} selectedKeys={[current]} mode="horizontal" items={menuItems} />
    )
};

export default Menus;