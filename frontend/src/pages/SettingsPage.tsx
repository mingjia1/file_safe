import { useEffect, useState } from 'react';
import { Card, Form, InputNumber, Switch, Button, message, Divider } from 'antd';
import { adminAPI } from '../api/client';

export default function SettingsPage() {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    try {
      const response = await adminAPI.getEncryptionConfig();
      form.setFieldsValue({
        aes_key_length: response.data.aes_key_length,
        rsa_key_length: response.data.rsa_key_length,
        password_min_length: response.data.password_min_length,
        password_require_special: response.data.password_require_special,
        password_require_uppercase: response.data.password_require_uppercase,
        password_require_lowercase: response.data.password_require_lowercase,
        password_require_digit: response.data.password_require_digit,
        config_encrypt: response.data.config_encrypt,
        enable_signature: response.data.enable_signature,
      });
    } catch (error) {
      message.error('Failed to load encryption config');
    }
  };

  const handleSave = async (values: any) => {
    setLoading(true);
    try {
      await adminAPI.updateEncryptionConfig(values);
      message.success('Settings saved successfully');
    } catch (error) {
      message.error('Failed to save settings');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: 24 }}>
      <Card title="Encryption Settings">
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSave}
          initialValues={{
            aes_key_length: 256,
            rsa_key_length: 2048,
            password_min_length: 8,
            password_require_special: true,
            password_require_uppercase: true,
            password_require_lowercase: true,
            password_require_digit: true,
            config_encrypt: true,
            enable_signature: true,
          }}
        >
          <Divider>Key Lengths</Divider>
          <Form.Item name="aes_key_length" label="AES Key Length">
            <InputNumber min={128} max={256} step={64} />
          </Form.Item>
          <Form.Item name="rsa_key_length" label="RSA Key Length">
            <InputNumber min={1024} max={4096} step={1024} />
          </Form.Item>

          <Divider>Password Policy</Divider>
          <Form.Item name="password_min_length" label="Minimum Password Length">
            <InputNumber min={4} max={64} />
          </Form.Item>
          <Form.Item name="password_require_special" label="Require Special Characters" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="password_require_uppercase" label="Require Uppercase" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="password_require_lowercase" label="Require Lowercase" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="password_require_digit" label="Require Digits" valuePropName="checked">
            <Switch />
          </Form.Item>

          <Divider>Security</Divider>
          <Form.Item name="config_encrypt" label="Encrypt Configuration" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="enable_signature" label="Enable Digital Signature" valuePropName="checked">
            <Switch />
          </Form.Item>

          <Button type="primary" htmlType="submit" loading={loading}>
            Save Settings
          </Button>
        </Form>
      </Card>
    </div>
  );
}
