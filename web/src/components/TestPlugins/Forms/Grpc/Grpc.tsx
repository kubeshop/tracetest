import {Form, Select, Tabs} from 'antd';
import {useEffect, useState} from 'react';
import GrpcService from 'services/Triggers/Grpc.service';
import {SupportedEditors} from 'constants/Editor.constants';
import {Editor, FileUpload} from 'components/Inputs';
import {Auth, Metadata, SkipTraceCollection} from 'components/Fields';
import TriggerTab from 'components/TriggerTab';
import useQueryTabs from 'hooks/useQueryTabs';
import {TDraftTest} from 'types/Test.types';

const RequestDetailsForm = () => {
  const [methodList, setMethodList] = useState<string[]>([]);
  const form = Form.useFormInstance<TDraftTest>();
  const protoFile = Form.useWatch('protoFile', form);
  const authType = Form.useWatch(['auth', 'type'], form);
  const message = Form.useWatch('message', form);
  const metadata = Form.useWatch('metadata', form);
  const skipTraceCollection = Form.useWatch('skipTraceCollection', form);

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

  const [activeKey, setActiveKey] = useQueryTabs('service-definition', 'triggerTab');

  return (
    <Tabs defaultActiveKey={activeKey} onChange={setActiveKey} activeKey={activeKey}>
      <Tabs.TabPane
        forceRender
        tab={<TriggerTab hasContent={!!protoFile} label="Service definition" />}
        key="service-definition"
      >
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

      <Tabs.TabPane forceRender tab={<TriggerTab hasContent={!!authType} label="Auth" />} key="auth">
        <Auth />
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab={<TriggerTab hasContent={!!message} label="Message" />} key="message">
        <Form.Item data-cy="message" name="message" style={{marginBottom: 0}}>
          <Editor
            type={SupportedEditors.Interpolation}
            placeholder="Enter message"
            basicSetup={{lineNumbers: true}}
            indentWithTab
          />
        </Form.Item>
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab={<TriggerTab totalItems={metadata?.length} label="Metadata" />} key="metadata">
        <Metadata />
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab={<TriggerTab hasContent={!!skipTraceCollection} label="Settings" />} key="settings">
        <SkipTraceCollection />
      </Tabs.TabPane>
    </Tabs>
  );
};

export default RequestDetailsForm;
