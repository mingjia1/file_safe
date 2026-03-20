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
      console.error('Failed to load audit logs');
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
      console.error('Failed to export');
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 220 },
    { title: 'Action', dataIndex: 'action', key: 'action', render: (action: string) => (
      <span style={{ color: action.includes('SUCCESS') ? 'green' : action.includes('FAIL') ? 'red' : 'blue' }}>
        {action}
      </span>
    )},
    { title: 'Package', dataIndex: 'package_name', key: 'package_name' },
    { title: 'User', dataIndex: 'username', key: 'username' },
    { title: 'IP Address', dataIndex: 'ip_address', key: 'ip_address' },
    { title: 'Detail', dataIndex: 'detail', key: 'detail', render: (detail: any) => detail ? JSON.stringify(detail) : '-' },
    { title: 'Created', dataIndex: 'created_at', key: 'created_at', render: (date: string) => new Date(date).toLocaleString() },
  ];

  return (
    <div style={{ padding: 24 }}>
      <Card 
        title="Audit Logs"
        extra={
          <Space>
            <Select allowClear placeholder="Filter by action" style={{ width: 160 }} onChange={setAction}>
              <Select.Option value="VERIFY_SUCCESS">Verify Success</Select.Option>
              <Select.Option value="VERIFY_FAIL">Verify Fail</Select.Option>
              <Select.Option value="PACKAGE_CREATE">Package Create</Select.Option>
              <Select.Option value="PACKAGE_DELETE">Package Delete</Select.Option>
              <Select.Option value="POLICY_CREATE">Policy Create</Select.Option>
              <Select.Option value="DOWNLOAD">Download</Select.Option>
            </Select>
            <Button icon={<ExportOutlined />} onClick={handleExport}>Export</Button>
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
            showTotal: (total) => `Total ${total} items`,
          }}
        />
      </Card>
    </div>
  );
}
