import {Col, Form, Row} from 'antd';
import TracetestAPI from 'redux/apis/Tracetest';
import {TestSelection} from 'components/Inputs';

const {useGetTestListQuery} = TracetestAPI.instance;

const TestsSelectionForm = () => {
  const {data} = useGetTestListQuery({take: 1000, skip: 0});

  return (
    <Row gutter={12}>
      <Col span={18}>
        <Form.Item name="steps">
          <TestSelection testList={data?.items || []} />
        </Form.Item>
      </Col>
    </Row>
  );
};

export default TestsSelectionForm;
