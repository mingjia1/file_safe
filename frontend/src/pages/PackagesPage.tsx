import { useEffect, useState } from 'react';
import { Table, Card, Button, Space, Modal, Form, Input, Upload, Select, message, Popconfirm } from 'antd';
import { PlusOutlined, DeleteOutlined, DownloadOutlined, KeyOutlined } from '@ant-design/icons';
import { packageAPI, passwordAPI } from '../api/client';

const { Option } = Select;

export default function PackagesPage() {
  const [packages, setPackages] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalVisible, setModalVisible] = useState(false);
  const [passwordModalVisible, setPasswordModalVisible] = useState(false);
  const [selectedPackage, setSelectedPackage] = useState<any>(null);
  const [passwords, setPasswords] = useState<any[]>([]);
  const [form] = Form.useForm();
  const [passwordForm] = Form.useForm();

  useEffect(() => {
    loadPackages();
  }, []);

  const loadPackages = async () => {
    try {
      const response = await packageAPI.list();
      setPackages(response.data.items);
    } catch (error) {
      message.error('加载文件包失败');
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (values: any) => {
    try {
      const formData = new FormData();
      formData.append('name', values.name);
      formData.append('format', values.format);
      formData.append('description', values.description || '');
      formData.append('file', values.file.file);
      await packageAPI.create(formData);
      message.success('文件包创建成功');
      setModalVisible(false);
      form.resetFields();
      loadPackages();
    } catch (error: any) {
      message.error(error.response?.data?.detail || '创建文件包失败');
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await packageAPI.delete(id);
      message.success('文件包删除成功');
      loadPackages();
    } catch (error) {
      message.error('删除文件包失败');
    }
  };

  const handleDownload = async (id: string, name: string) => {
    try {
      const response = await packageAPI.download(id);
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', name);
      document.body.appendChild(link);
      link.click();
      link.remove();
    } catch (error) {
      message.error('下载文件包失败');
    }
  };

  const handleAddPassword = async (values: { password: string; priority: number }) => {
    if (!selectedPackage) return;
    try {
      await passwordAPI.create(selectedPackage.id, values);
      message.success('密码添加成功');
      setPasswordModalVisible(false);
      passwordForm.resetFields();
      loadPasswords(selectedPackage.id);
    } catch (error) {
      message.error('添加密码失败');
    }
  };

  const loadPasswords = async (packageId: string) => {
    try {
      const response = await passwordAPI.list(packageId);
      setPasswords(response.data);
    } catch (error) {
      message.error('加载密码失败');
    }
  };

  const openPasswordModal = async (pkg: any) => {
    setSelectedPackage(pkg);
    await loadPasswords(pkg.id);
    setPasswordModalVisible(true);
  };

  const columns = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '格式', dataIndex: 'format', key: 'format' },
    { title: '状态', dataIndex: 'status', key: 'status', render: (status: string) => (
      <span style={{ color: status === 'active' ? 'green' : 'gray' }}>{status}</span>
    )},
    { title: '文件大小', dataIndex: 'file_size', key: 'file_size', render: (size: number) => `${(size / 1024).toFixed(2)} KB` },
    { title: '密码数', dataIndex: 'password_count', key: 'password_count' },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', render: (date: string) => new Date(date).toLocaleString() },
    { title: '操作', key: 'actions', render: (_: any, record: any) => (
      <Space>
        <Button size="small" icon={<KeyOutlined />} onClick={() => openPasswordModal(record)}>密码</Button>
        <Button size="small" icon={<DownloadOutlined />} onClick={() => handleDownload(record.id, record.name)}>下载</Button>
        <Popconfirm title="确定删除?" onConfirm={() => handleDelete(record.id)}>
          <Button size="small" danger icon={<DeleteOutlined />}>删除</Button>
        </Popconfirm>
      </Space>
    )},
  ];

  const passwordColumns = [
    { title: '优先级', dataIndex: 'priority', key: 'priority' },
    { title: '状态', dataIndex: 'status', key: 'status' },
    { title: '生效时间', dataIndex: 'valid_from', key: 'valid_from', render: (d: string) => d ? new Date(d).toLocaleString() : '-' },
    { title: '失效时间', dataIndex: 'valid_until', key: 'valid_until', render: (d: string) => d ? new Date(d).toLocaleString() : '-' },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', render: (date: string) => new Date(date).toLocaleString() },
  ];

  return (
    <div style={{ padding: 24 }}>
      <Card 
        title="文件包管理" 
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
            创建文件包
          </Button>
        }
      >
        <Table dataSource={packages} columns={columns} rowKey="id" loading={loading} />
      </Card>

      <Modal title="创建文件包" open={modalVisible} onCancel={() => setModalVisible(false)} footer={null}>
        <Form form={form} layout="vertical" onFinish={handleCreate}>
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="format" label="格式" rules={[{ required: true }]}>
            <Select>
              <Option value="exe">EXE</Option>
              <Option value="zip">ZIP</Option>
            </Select>
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea />
          </Form.Item>
          <Form.Item name="file" label="文件" rules={[{ required: true }]}>
            <Upload beforeUpload={() => false} maxCount={1}>
              <Button>选择文件</Button>
            </Upload>
          </Form.Item>
          <Button type="primary" htmlType="submit">创建</Button>
        </Form>
      </Modal>

      <Modal 
        title={`${selectedPackage?.name} - 密码管理`} 
        open={passwordModalVisible} 
        onCancel={() => setPasswordModalVisible(false)}
        footer={null}
        width={700}
      >
        <div style={{ marginBottom: 16 }}>
          <Form form={passwordForm} layout="inline" onFinish={handleAddPassword}>
            <Form.Item name="password" rules={[{ required: true }]}>
              <Input placeholder="密码" />
            </Form.Item>
            <Form.Item name="priority" initialValue={1}>
              <Input type="number" placeholder="优先级" />
            </Form.Item>
            <Button type="primary" htmlType="submit">添加密码</Button>
          </Form>
        </div>
        <Table dataSource={passwords} columns={passwordColumns} rowKey="id" size="small" />
      </Modal>
    </div>
  );
}
