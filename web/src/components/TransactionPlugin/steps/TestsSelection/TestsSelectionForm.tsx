import {Col, Form, Row} from 'antd';
import {useGetTestListQuery} from 'redux/apis/TraceTest.api';
import TestsSelectionInput from './TestsSelectionInput/TestsSelectionInput';

const TestsSelectionForm = () => {
  const {data} = useGetTestListQuery({take: 1000, skip: 0});

  return (
    <Row gutter={12}>
      <Col span={18}>
        <Form.Item name="tests">
          <TestsSelectionInput testList={data?.items || []} />
        </Form.Item>
      </Col>
    </Row>
  );
};

export default TestsSelectionForm;
