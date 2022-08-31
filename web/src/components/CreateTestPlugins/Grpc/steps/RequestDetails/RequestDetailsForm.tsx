import {Col, Form, Input, Row, Select} from 'antd';
import {useCallback, useEffect, useState} from 'react';
import GrpcService from 'services/Triggers/Grpc.service';
import {IRpcValues, TDraftTestForm} from 'types/Test.types';
import RequestDetailsAuthInput from '../../../Rest/steps/RequestDetails/RequestDetailsAuthInput/RequestDetailsAuthInput';
import RequestDetailsUrlInput from '../../../Rest/steps/RequestDetails/RequestDetailsUrlInput';
import RequestDetailsFileInput from './RequestDetailsFileInput';
import RequestDetailsMetadataInput from './RequestDetailsMetadataInput';

interface IProps {
  form: TDraftTestForm<IRpcValues>;
}

const RequestDetailsForm = ({form}: IProps) => {
  const [methodList, setMethodList] = useState<string[]>([]);
  const protoFile = Form.useWatch('protoFile', form);

  const getMethodList = useCallback(async () => {
    if (protoFile) {
      const fileText = await protoFile.text();
      const list = GrpcService.getMethodList(fileText);

      setMethodList(list);
    } else {
      setMethodList([]);
    }
  }, [protoFile]);

  useEffect(() => {
    getMethodList();
  }, [getMethodList]);

  return (
    <>
      <Row gutter={12}>
        <Col span={12}>
          <span>
            <Form.Item data-cy="protoFile" name="protoFile" label="Upload Protobuf File">
              <RequestDetailsFileInput />
            </Form.Item>
          </span>
        </Col>
        <Col span={12}>
          <Form.Item data-cy="method" label="Select Method" name="method">
            <Select data-cy="method-select">
              {methodList.map(method => (
                <Select.Option data-cy={`rpc-method-${method}`} key={method} value={method}>
                  {method}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
        </Col>
      </Row>
      <Row gutter={12}>
        <Col span={24}>
          <RequestDetailsUrlInput showMethodSelector={false} shouldValidateUrl={false} />
        </Col>
      </Row>
      <Row gutter={12}>
        <Col span={24}>
          <RequestDetailsAuthInput form={form} />
        </Col>
      </Row>
      <Row gutter={12}>
        <Col span={12}>
          <RequestDetailsMetadataInput />
        </Col>
      </Row>
      <Row gutter={12}>
        <Col span={24}>
          <Form.Item data-cy="message" label="Message" name="message" style={{marginBottom: 0}}>
            <Input.TextArea placeholder="Enter message" />
          </Form.Item>
        </Col>
      </Row>
    </>
  );
};

export default RequestDetailsForm;
