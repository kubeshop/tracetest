import {Col, Form, Row, Select} from 'antd';
import {useCallback, useEffect, useState} from 'react';
import GrpcService from 'services/Triggers/Grpc.service';
import {IRpcValues, TDraftTestForm} from 'types/Test.types';
import {SupportedEditors} from 'constants/Editor.constants';
import Editor from 'components/Editor';
import RequestDetailsAuthInput from '../../../Rest/steps/RequestDetails/RequestDetailsAuthInput/RequestDetailsAuthInput';
import RequestDetailsUrlInput from '../../../Rest/steps/RequestDetails/RequestDetailsUrlInput';
import RequestDetailsFileInput from './RequestDetailsFileInput';
import RequestDetailsMetadataInput from './RequestDetailsMetadataInput';
import * as S from './RequestDetails.styled';

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
    <S.FieldsContainer>
      <Row gutter={12}>
        <Col span={18}>
          <span>
            <Form.Item data-cy="protoFile" name="protoFile" label="Upload Protobuf File">
              <RequestDetailsFileInput />
            </Form.Item>
          </span>
        </Col>
      </Row>
      <Row gutter={12}>
        <Col span={18}>
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
        <Col span={18}>
          <RequestDetailsUrlInput showMethodSelector={false} />
        </Col>
      </Row>
      <Row gutter={12}>
        <Col span={18}>
          <RequestDetailsAuthInput />
        </Col>
      </Row>
      <Row gutter={12}>
        <Col span={12}>
          <RequestDetailsMetadataInput />
        </Col>
      </Row>
      <Row gutter={12}>
        <Col span={18}>
          <Form.Item data-cy="message" label="Message" name="message" style={{marginBottom: 0}}>
            <Editor
              type={SupportedEditors.Interpolation}
              placeholder="Enter message"
              basicSetup={{lineNumbers: true}}
              indentWithTab
            />
          </Form.Item>
        </Col>
      </Row>
    </S.FieldsContainer>
  );
};

export default RequestDetailsForm;
