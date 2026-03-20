import { useEffect, useState } from 'react';
import { Table, Card, Select, Space, Button } from 'antd';
import { ExportOutlined } from '@ant-design/icons';
import { auditAPI } from '../api/client';

export default function AuditPage() {
  const [logs, setLogs] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [action, setAction] = useState<string>('');

  useEffect(() => {
    loadLogs();
  }, [page, action]);

  const loadLogs = async () => {
    setLoading(true);
    try {
      const params: any = { page, page_size: 20 };
      if (action) params.action = action;
      const response = await auditAPI.list(params);
      setLogs(response.data.items);
      setTotal(response.data.total);
    } catch (error) {
      console.error('加载审计日志失败');
    } finally {
      setLoading(false);
    }
  };

  const handleExport = async () => {
    try {
      const response = await auditAPI.export('csv');
      const blob = new Blob([response.data.content], { type: 'text/csv' });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `audit_logs_${Date.now()}.csv`);
      document.body.appendChild(link);
      link.click();
      link.remove();
    } catch (error) {
      console.error('导出失败');
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 220 },
    { title: '操作', dataIndex: 'action', key: 'action', render: (action: string) => (
      <span style={{ color: action.includes('SUCCESS') ? 'green' : action.includes('FAIL') ? 'red' : 'blue' }}>
        {action}
      </span>
    )},
    { title: '文件包', dataIndex: 'package_name', key: 'package_name' },
    { title: '用户', dataIndex: 'username', key: 'username' },
    { title: 'IP 地址', dataIndex: 'ip_address', key: 'ip_address' },
    { title: '详情', dataIndex: 'detail', key: 'detail', render: (detail: any) => detail ? JSON.stringify(detail) : '-' },
    { title: '时间', dataIndex: 'created_at', key: 'created_at', render: (date: string) => new Date(date).toLocaleString() },
  ];

  return (
    <div style={{ padding: 24 }}>
      <Card 
        title="审计日志"
        extra={
          <Space>
            <Select allowClear placeholder="按操作筛选" style={{ width: 160 }} onChange={setAction}>
              <Select.Option value="VERIFY_SUCCESS">验证成功</Select.Option>
              <Select.Option value="VERIFY_FAIL">验证失败</Select.Option>
              <Select.Option value="PACKAGE_CREATE">创建文件包</Select.Option>
              <Select.Option value="PACKAGE_DELETE">删除文件包</Select.Option>
              <Select.Option value="POLICY_CREATE">创建策略</Select.Option>
              <Select.Option value="DOWNLOAD">下载</Select.Option>
            </Select>
            <Button icon={<ExportOutlined />} onClick={handleExport}>导出</Button>
          </Space>
        }
      >
        <Table 
          dataSource={logs} 
          columns={columns} 
          rowKey="id" 
          loading={loading}
          pagination={{
            current: page,
            total,
            pageSize: 20,
            onChange: setPage,
            showTotal: (total) => `共 ${total} 条`,
          }}
        />
      </Card>
    </div>
  );
}
