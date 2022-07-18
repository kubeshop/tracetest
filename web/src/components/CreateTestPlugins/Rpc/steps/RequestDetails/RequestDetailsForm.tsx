import {Form, Input, Select} from 'antd';
import {useCallback, useState, useEffect} from 'react';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import {IRpcValues, TDraftTestForm} from 'types/Test.types';
import RpcService from 'services/Triggers/Rpc.service';
import RequestDetailsAuthInput from '../../../Rest/steps/RequestDetails/RequestDetailsAuthInput/RequestDetailsAuthInput';
import RequestDetailsUrlInput from '../../../Rest/steps/RequestDetails/RequestDetailsUrlInput';
import * as S from './RequestDetails.styled';
import RequestDetailsMetadataInput from './RequestDetailsMetadataInput';
import RequestDetailsFileInput from './RequestDetailsFileInput';

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
    <>
      <S.InputContainer>
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
      </S.InputContainer>
      <Step.Title>Provide Additional Information</Step.Title>
      <S.DoubleInputContainer>
        <RequestDetailsUrlInput showMethodSelector={false} shouldValidateUrl={false} />
        <RequestDetailsAuthInput form={form} />
        <RequestDetailsMetadataInput />
        <Form.Item data-cy="message" label="Message" name="message" style={{marginBottom: 0}}>
          <Input.TextArea placeholder="Enter message" />
        </Form.Item>
      </S.DoubleInputContainer>
    </>
  );
};

export default RequestDetailsForm;
