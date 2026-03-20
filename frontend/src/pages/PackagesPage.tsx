import { useEffect, useState } from 'react';
import { Table, Card, Button, Space, Modal, Form, Input, InputNumber, Upload, Select, message, Popconfirm, DatePicker, Tag } from 'antd';
import { PlusOutlined, DeleteOutlined, DownloadOutlined, KeyOutlined, EditOutlined, CheckCircleOutlined, StopOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import { packageAPI, passwordAPI } from '../api/client';

const { Option } = Select;

export default function PackagesPage() {
  const [packages, setPackages] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [passwordModalVisible, setPasswordModalVisible] = useState(false);
  const [editingPackage, setEditingPackage] = useState<any>(null);
  const [editingPassword, setEditingPassword] = useState<any>(null);
  const [selectedPackage, setSelectedPackage] = useState<any>(null);
  const [passwords, setPasswords] = useState<any[]>([]);
  const [form] = Form.useForm();
  const [editForm] = Form.useForm();
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
      if (values.file) {
        formData.append('file', values.file.file);
      }
      await packageAPI.create(formData);
      message.success('文件包创建成功');
      setCreateModalVisible(false);
      form.resetFields();
      loadPackages();
    } catch (error: any) {
      message.error(error.response?.data?.detail || '创建文件包失败');
    }
  };

  const handleEdit = async (values: any) => {
    if (!editingPackage) return;
    try {
      await packageAPI.update(editingPackage.id, {
        name: values.name,
        description: values.description,
        status: values.status,
      });
      message.success('文件包更新成功');
      setEditModalVisible(false);
      editForm.resetFields();
      setEditingPackage(null);
      loadPackages();
    } catch (error: any) {
      message.error(error.response?.data?.detail || '更新文件包失败');
    }
  };

  const openEditModal = (pkg: any) => {
    setEditingPackage(pkg);
    editForm.setFieldsValue({
      name: pkg.name,
      description: pkg.description,
      status: pkg.status,
    });
    setEditModalVisible(true);
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

  const handleAddPassword = async (values: any) => {
    if (!selectedPackage) return;
    if (!values.password) {
      message.error('请输入密码');
      return;
    }
    try {
      const passwordData: any = {
        password: values.password,
        priority: values.priority || 1,
      };
      if (values.valid_from) {
        passwordData.valid_from = values.valid_from.toISOString();
      }
      if (values.valid_until) {
        passwordData.valid_until = values.valid_until.toISOString();
      }
      await passwordAPI.create(selectedPackage.id, passwordData);
      message.success('密码添加成功');
      passwordForm.resetFields();
      loadPasswords(selectedPackage.id);
    } catch (error: any) {
      message.error(error.response?.data?.detail || '添加密码失败');
    }
  };

  const handleUpdatePassword = async (values: any) => {
    if (!editingPassword) return;
    try {
      const passwordData: any = {};
      if (values.password) passwordData.password = values.password;
      if (values.priority) passwordData.priority = values.priority;
      if (values.valid_from) passwordData.valid_from = values.valid_from.toISOString();
      if (values.valid_until) passwordData.valid_until = values.valid_until.toISOString();
      
      await passwordAPI.update(editingPassword.id, passwordData);
      message.success('密码更新成功');
      setEditingPassword(null);
      passwordForm.resetFields();
      loadPasswords(selectedPackage.id);
    } catch (error: any) {
      message.error(error.response?.data?.detail || '更新密码失败');
    }
  };

  const handleDeletePassword = async (id: string) => {
    try {
      await passwordAPI.delete(id);
      message.success('密码删除成功');
      loadPasswords(selectedPackage.id);
    } catch (error) {
      message.error('删除密码失败');
    }
  };

  const handleActivatePassword = async (id: string) => {
    try {
      await passwordAPI.activate(id);
      message.success('密码已激活');
      loadPasswords(selectedPackage.id);
    } catch (error) {
      message.error('激活密码失败');
    }
  };

  const handleDeactivatePassword = async (id: string) => {
    try {
      await passwordAPI.deactivate(id);
      message.success('密码已停用');
      loadPasswords(selectedPackage.id);
    } catch (error) {
      message.error('停用密码失败');
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

  const openEditPasswordModal = (pwd: any) => {
    setEditingPassword(pwd);
    passwordForm.setFieldsValue({
      password: '',
      priority: pwd.priority,
      valid_from: pwd.valid_from ? dayjs(pwd.valid_from) : null,
      valid_until: pwd.valid_until ? dayjs(pwd.valid_until) : null,
    });
  };

  const columns = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '格式', dataIndex: 'format', key: 'format' },
    { title: '状态', dataIndex: 'status', key: 'status', render: (status: string) => (
      <Tag color={status === 'active' ? 'green' : 'red'}>{status}</Tag>
    )},
    { title: '文件大小', dataIndex: 'file_size', key: 'file_size', render: (size: number) => `${(size / 1024).toFixed(2)} KB` },
    { title: '密码数', dataIndex: 'password_count', key: 'password_count' },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', render: (date: string) => new Date(date).toLocaleString() },
    { title: '操作', key: 'actions', width: 280, render: (_: any, record: any) => (
      <Space size="small">
        <Button size="small" type="link" icon={<EditOutlined />} onClick={() => openEditModal(record)}>编辑</Button>
        <Button size="small" icon={<KeyOutlined />} onClick={() => openPasswordModal(record)}>密码</Button>
        <Button size="small" icon={<DownloadOutlined />} onClick={() => handleDownload(record.id, record.name)}>下载</Button>
        <Popconfirm title="确定删除?" onConfirm={() => handleDelete(record.id)}>
          <Button size="small" danger type="link" icon={<DeleteOutlined />}>删除</Button>
        </Popconfirm>
      </Space>
    )},
  ];

  const passwordColumns = [
    { title: '优先级', dataIndex: 'priority', key: 'priority', width: 80 },
    { title: '状态', dataIndex: 'status', key: 'status', width: 100, render: (status: string) => (
      <Tag color={status === 'active' ? 'green' : 'red'}>{status}</Tag>
    )},
    { title: '生效时间', dataIndex: 'valid_from', key: 'valid_from', width: 180, render: (d: string) => d ? new Date(d).toLocaleString() : '-' },
    { title: '失效时间', dataIndex: 'valid_until', key: 'valid_until', width: 180, render: (d: string) => d ? new Date(d).toLocaleString() : '-' },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 180, render: (date: string) => new Date(date).toLocaleString() },
    { title: '操作', key: 'actions', width: 200, render: (_: any, record: any) => (
      <Space size="small">
        <Button size="small" type="link" icon={<EditOutlined />} onClick={() => openEditPasswordModal(record)}>编辑</Button>
        {record.status === 'active' ? (
          <Button size="small" type="link" icon={<StopOutlined />} onClick={() => handleDeactivatePassword(record.id)}>停用</Button>
        ) : (
          <Button size="small" type="link" icon={<CheckCircleOutlined />} onClick={() => handleActivatePassword(record.id)}>激活</Button>
        )}
        <Popconfirm title="确定删除?" onConfirm={() => handleDeletePassword(record.id)}>
          <Button size="small" danger type="link" icon={<DeleteOutlined />}>删除</Button>
        </Popconfirm>
      </Space>
    )},
  ];

  return (
    <div style={{ padding: 24 }}>
      <Card 
        title="文件包管理" 
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>
            创建文件包
          </Button>
        }
      >
        <Table dataSource={packages} columns={columns} rowKey="id" loading={loading} />
      </Card>

      {/* 创建文件包 */}
      <Modal 
        title="创建文件包" 
        open={createModalVisible} 
        onCancel={() => { setCreateModalVisible(false); form.resetFields(); }}
        footer={null}
      >
        <Form form={form} layout="vertical" onFinish={handleCreate}>
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input placeholder="输入文件包名称" />
          </Form.Item>
          <Form.Item name="format" label="格式" rules={[{ required: true }]}>
            <Select placeholder="选择格式">
              <Option value="exe">EXE</Option>
              <Option value="zip">ZIP</Option>
            </Select>
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} placeholder="输入描述" />
          </Form.Item>
          <Form.Item name="file" label="文件" rules={[{ required: true }]}>
            <Upload beforeUpload={() => false} maxCount={1}>
              <Button>选择文件</Button>
            </Upload>
          </Form.Item>
          <Button type="primary" htmlType="submit" block>创建</Button>
        </Form>
      </Modal>

      {/* 编辑文件包 */}
      <Modal 
        title="编辑文件包" 
        open={editModalVisible} 
        onCancel={() => { setEditModalVisible(false); editForm.resetFields(); setEditingPackage(null); }}
        footer={null}
      >
        <Form form={editForm} layout="vertical" onFinish={handleEdit}>
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item name="status" label="状态" rules={[{ required: true }]}>
            <Select>
              <Option value="active">激活</Option>
              <Option value="inactive">停用</Option>
            </Select>
          </Form.Item>
          <Space>
            <Button type="primary" htmlType="submit">保存</Button>
            <Button onClick={() => { setEditModalVisible(false); editForm.resetFields(); }}>取消</Button>
          </Space>
        </Form>
      </Modal>

      {/* 密码管理 */}
      <Modal 
        title={`${selectedPackage?.name} - 密码管理`} 
        open={passwordModalVisible} 
        onCancel={() => { setPasswordModalVisible(false); setSelectedPackage(null); passwordForm.resetFields(); }}
        footer={null}
        width={900}
      >
        {/* 添加/编辑密码表单 */}
        <Card size="small" style={{ marginBottom: 16 }}>
          <Form 
            form={passwordForm} 
            layout="inline" 
            onFinish={editingPassword ? handleUpdatePassword : handleAddPassword}
          >
            <Space align="start" wrap>
              <Form.Item name="password" label="密码" rules={editingPassword ? [] : [{ required: true }]}>
                <Input.Password placeholder={editingPassword ? "留空则不修改" : "输入密码"} style={{ width: 160 }} />
              </Form.Item>
              <Form.Item name="priority" label="优先级" initialValue={editingPassword?.priority || 1}>
                <InputNumber min={1} style={{ width: 80 }} />
              </Form.Item>
              <Form.Item name="valid_from" label="生效时间">
                <DatePicker showTime placeholder="生效时间" />
              </Form.Item>
              <Form.Item name="valid_until" label="失效时间">
                <DatePicker showTime placeholder="失效时间" />
              </Form.Item>
              <Form.Item>
                <Space>
                  <Button type="primary" htmlType="submit">{editingPassword ? '更新' : '添加'}</Button>
                  {editingPassword && (
                    <Button onClick={() => { setEditingPassword(null); passwordForm.resetFields(); }}>取消</Button>
                  )}
                </Space>
              </Form.Item>
            </Space>
          </Form>
        </Card>
        
        <Table dataSource={passwords} columns={passwordColumns} rowKey="id" size="small" pagination={false} />
      </Modal>
    </div>
  );
}
