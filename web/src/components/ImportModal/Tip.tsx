import {BulbOutlined} from '@ant-design/icons';
import {Typography} from 'antd';

const Tip = () => {
  return (
    <>
      <Typography.Title level={3} type="secondary">
        <BulbOutlined /> What are supported formats?
      </Typography.Title>
      <Typography.Text type="secondary">We support cURL & Postman. OpenAPI is coming soon!</Typography.Text>
    </>
  );
};

export default Tip;
