import {Form, Select, Tabs} from 'antd';
import {useEffect, useState} from 'react';
import GrpcService from 'services/Triggers/Grpc.service';
import {IRpcValues, TDraftTestForm} from 'types/Test.types';
import {SupportedEditors} from 'constants/Editor.constants';
import Editor from 'components/Editor';
import RequestDetailsFileInput from './RequestDetailsFileInput';
import RequestDetailsMetadataInput from './RequestDetailsMetadataInput';
import RequestDetailsAuthInput from '../../../Rest/steps/RequestDetails/RequestDetailsAuthInput/RequestDetailsAuthInput';

interface IProps {
  form: TDraftTestForm<IRpcValues>;
}

const RequestDetailsForm = ({form}: IProps) => {
  const [methodList, setMethodList] = useState<string[]>([]);
  const protoFile = Form.useWatch('protoFile', form);

  useEffect(() => {
    const getMethodList = async () => {
      if (protoFile) {
        const fileText = await protoFile.text();
        const list = GrpcService.getMethodList(fileText);

        setMethodList(list);
      } else {
        setMethodList([]);
      }
    };

    getMethodList();
  }, [protoFile]);

  return (
    <Tabs defaultActiveKey="general">
      <Tabs.TabPane forceRender tab="General" key="general">
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

        <RequestDetailsAuthInput />
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab="Message" key="message">
        <Form.Item data-cy="message" label="Message" name="message" style={{marginBottom: 0}}>
          <Editor
            type={SupportedEditors.Interpolation}
            placeholder="Enter message"
            basicSetup={{lineNumbers: true}}
            indentWithTab
          />
        </Form.Item>
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab="Metadata" key="metadata">
        <RequestDetailsMetadataInput />
      </Tabs.TabPane>
    </Tabs>
  );
};

export default RequestDetailsForm;
