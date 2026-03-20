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
      console.error('加载统计失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: 24 }}>
      <h1>仪表盘</h1>
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="文件包总数"
              value={stats?.total_packages || 0}
              prefix={<FileProtectOutlined />}
              loading={loading}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="活跃文件包"
              value={stats?.active_packages || 0}
              prefix={<FileProtectOutlined />}
              loading={loading}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="活跃密码"
              value={stats?.active_passwords || 0}
              prefix={<LockOutlined />}
              loading={loading}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="验证成功率"
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
              title="总下载次数"
              value={stats?.total_downloads || 0}
              prefix={<DownloadOutlined />}
              loading={loading}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="总验证次数"
              value={stats?.total_verifies || 0}
              prefix={<AuditOutlined />}
              loading={loading}
            />
          </Card>
        </Col>
      </Row>

      <Card 
        title="快捷操作"
        extra={
          <Space>
            <Button type="primary" onClick={() => navigate('/packages')}>
              管理文件包
            </Button>
            <Button onClick={() => navigate('/audit')}>
              查看审计日志
            </Button>
          </Space>
        }
      >
        <p>欢迎使用密码时效管理器。使用导航菜单管理文件包、密码策略和查看审计日志。</p>
      </Card>
    </div>
  );
}
