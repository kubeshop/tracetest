import {Col, Form, Input, Row, Select} from 'antd';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import {useCallback, useEffect, useState} from 'react';
import RpcService from 'services/Triggers/Rpc.service';
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
      const list = RpcService.getMethodList(fileText);

      setMethodList(list);
    } else {
      setMethodList([]);
      form.setFieldsValue({
        method: '',
      });
    }
  }, [form, protoFile]);

  useEffect(() => {
    getMethodList();
  }, [getMethodList]);

  return (
    <div style={{display: 'grid'}}>
      <Row gutter={12}>
        <Col span={12}>
          <span>
            <Form.Item data-cy="protoFile" name="protoFile" label="Upload Protobuf File">
              <RequestDetailsFileInput />
            </Form.Item>
            <Form.Item data-cy="method" label="Select Method" name="method">
              <Select data-cy="method-select">
                {methodList.map(method => (
                  <Select.Option data-cy={`rpc-method-${method}`} key={method} value={method}>
                    {method}
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
          </span>
        </Col>
      </Row>
      <Step.Title>Provide Additional Information</Step.Title>
      <Row gutter={12}>
        <Col span={12}>
          <RequestDetailsUrlInput showMethodSelector={false} shouldValidateUrl={false} />
        </Col>
        <Col span={12}>
          <RequestDetailsAuthInput form={form} />
        </Col>
      </Row>
      <Row gutter={12} style={{marginTop: 16}}>
        <Col span={12}>
          <RequestDetailsMetadataInput />
        </Col>
        <Col span={12}>
          <Form.Item data-cy="message" label="Message" name="message" style={{marginBottom: 0}}>
            <Input.TextArea placeholder="Enter message" />
          </Form.Item>
        </Col>
      </Row>
    </div>
  );
};

export default RequestDetailsForm;
