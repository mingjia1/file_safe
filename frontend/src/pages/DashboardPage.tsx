import { useEffect, useState } from 'react';
import { Card, Row, Col, Statistic, Space, Button } from 'antd';
import { FileProtectOutlined, LockOutlined, CheckCircleOutlined, DownloadOutlined, AuditOutlined } from '@ant-design/icons';
import { adminAPI } from '../api/client';
import { useNavigate } from 'react-router-dom';

export default function DashboardPage() {
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      const response = await adminAPI.getDashboardStats();
      setStats(response.data);
    } catch (error) {
      console.error('Failed to load statistics');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: 24 }}>
      <h1>Dashboard</h1>
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="Total Packages"
              value={stats?.total_packages || 0}
              prefix={<FileProtectOutlined />}
              loading={loading}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="Active Packages"
              value={stats?.active_packages || 0}
              prefix={<FileProtectOutlined />}
              loading={loading}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="Active Passwords"
              value={stats?.active_passwords || 0}
              prefix={<LockOutlined />}
              loading={loading}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="Verify Success Rate"
              value={(stats?.verify_success_rate || 0) * 100}
              suffix="%"
              prefix={<CheckCircleOutlined />}
              loading={loading}
              valueStyle={{ color: stats?.verify_success_rate > 0.5 ? '#3f8600' : '#cf1322' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="Total Downloads"
              value={stats?.total_downloads || 0}
              prefix={<DownloadOutlined />}
              loading={loading}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="Total Verifications"
              value={stats?.total_verifies || 0}
              prefix={<AuditOutlined />}
              loading={loading}
            />
          </Card>
        </Col>
      </Row>

      <Card 
        title="Quick Actions"
        extra={
          <Space>
            <Button type="primary" onClick={() => navigate('/packages')}>
              Manage Packages
            </Button>
            <Button onClick={() => navigate('/audit')}>
              View Audit Logs
            </Button>
          </Space>
        }
      >
        <p>Welcome to Password Timer Manager. Use the navigation menu to manage file packages, passwords, and view audit logs.</p>
      </Card>
    </div>
  );
}
