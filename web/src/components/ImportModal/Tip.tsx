import {Typography} from 'antd';

const Tip = () => (
  <>
    <Typography.Title level={3}>💡 What are the supported formats?</Typography.Title>
    <Typography.Text type="secondary">
      We support Tracetest Definition, cURL & Postman. OpenAPI is coming soon!
    </Typography.Text>
  </>
);

export default Tip;
