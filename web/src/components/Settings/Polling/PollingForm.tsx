import {useEffect} from 'react';
import {Col, Form, Input, Row} from 'antd';

import AllowButton, {Operation} from 'components/AllowButton';
import {useSettings} from 'providers/Settings/Settings.provider';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import {ResourceType, TDraftPollingProfiles} from 'types/Settings.types';
import SettingService from 'services/Setting.service';
import * as S from '../common/Settings.styled';

const FORM_ID = 'polling';

const PollingForm = () => {
  const [form] = Form.useForm<TDraftPollingProfiles>();
  const {isLoading, onSubmit} = useSettings();
  const {pollingProfile} = useSettingsValues();

  useEffect(() => {
    form.resetFields();
  }, [form, pollingProfile]);

  const handleOnSubmit = (values: TDraftPollingProfiles) => {
    onSubmit([SettingService.getDraftResource(ResourceType.PollingProfileType, values)]);
  };

  return (
    <Form<TDraftPollingProfiles>
      autoComplete="off"
      form={form}
      initialValues={pollingProfile}
      layout="vertical"
      name={FORM_ID}
      onFinish={handleOnSubmit}
    >
      <Form.Item hidden name="default" />
      <Form.Item hidden name="id" />
      <Form.Item hidden name="name" />
      <Form.Item hidden name="strategy" />

      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item
            label="Max wait time for trace"
            name={['periodic', 'timeout']}
            rules={[{required: true, message: 'Max wait time for trace is required'}]}
          >
            <Input placeholder="10s" />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item
            label="Retry delay"
            name={['periodic', 'retryDelay']}
            rules={[{required: true, message: 'Retry delay is required'}]}
          >
            <Input placeholder="500ms" />
          </Form.Item>
        </Col>
      </Row>

      <S.FooterContainer>
        <AllowButton operation={Operation.Configure} htmlType="submit" loading={isLoading} type="primary">
          Save
        </AllowButton>
      </S.FooterContainer>
    </Form>
  );
};

export default PollingForm;
