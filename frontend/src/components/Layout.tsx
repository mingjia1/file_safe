import { useState, type ReactNode } from 'react';
import { Layout, Menu } from 'antd';
import { DashboardOutlined, FileProtectOutlined, AuditOutlined, SettingOutlined, LogoutOutlined } from '@ant-design/icons';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuthStore } from '../stores/authStore';

const { Header, Content, Sider } = Layout;

interface AppLayoutProps {
  children: ReactNode;
}

export default function AppLayout({ children }: AppLayoutProps) {
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const logout = useAuthStore((state) => state.logout);

  const menuItems = [
    { key: '/dashboard', icon: <DashboardOutlined />, label: 'Dashboard' },
    { key: '/packages', icon: <FileProtectOutlined />, label: 'Packages' },
    { key: '/audit', icon: <AuditOutlined />, label: 'Audit Logs' },
    { key: '/settings', icon: <SettingOutlined />, label: 'Settings' },
  ];

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center', background: '#001529', padding: '0 24px' }}>
        <div style={{ color: 'white', fontSize: 18, fontWeight: 'bold', marginRight: 48 }}>
          Password Timer Manager
        </div>
        <div style={{ marginLeft: 'auto', color: 'white' }}>
          <LogoutOutlined style={{ cursor: 'pointer' }} onClick={logout} />
        </div>
      </Header>
      <Layout>
        <Sider collapsible collapsed={collapsed} onCollapse={setCollapsed}>
          <Menu 
            theme="dark" 
            mode="inline" 
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={({ key }) => navigate(key)}
          />
        </Sider>
        <Layout style={{ padding: '0 24px 24px' }}>
          <Content>
            {children}
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
}
