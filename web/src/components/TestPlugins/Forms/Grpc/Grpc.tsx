import {Form, Select, Tabs} from 'antd';
import {useEffect, useState} from 'react';
import GrpcService from 'services/Triggers/Grpc.service';
import {SupportedEditors} from 'constants/Editor.constants';
import {Editor, FileUpload} from 'components/Inputs';
import {Auth, Metadata, SkipTraceCollection} from 'components/Fields';
import useQueryTabs from 'hooks/useQueryTabs';

const RequestDetailsForm = () => {
  const [methodList, setMethodList] = useState<string[]>([]);
  const form = Form.useFormInstance();
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

  const [activeKey, setActiveKey] = useQueryTabs('service-definition');

  return (
    <Tabs defaultActiveKey={activeKey} onChange={setActiveKey} activeKey={activeKey}>
      <Tabs.TabPane forceRender tab="Service definition" key="service-definition">
        <Form.Item data-cy="protoFile" name="protoFile" label="Upload Protobuf File">
          <FileUpload />
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
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab="Auth" key="auth">
        <Auth />
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab="Message" key="message">
        <Form.Item data-cy="message" name="message" style={{marginBottom: 0}}>
          <Editor
            type={SupportedEditors.Interpolation}
            placeholder="Enter message"
            basicSetup={{lineNumbers: true}}
            indentWithTab
          />
        </Form.Item>
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab="Metadata" key="metadata">
        <Metadata />
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab="Settings" key="settings">
        <SkipTraceCollection />
      </Tabs.TabPane>
    </Tabs>
  );
};

export default RequestDetailsForm;
