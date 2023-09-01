import {Col, Form, Input, Row, Switch} from 'antd';
import {useEffect} from 'react';
import AllowButton, {Operation} from 'components/AllowButton';
import {useSettings} from 'providers/Settings/Settings.provider';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import SettingService from 'services/Setting.service';
import Collapse from 'components/Collapse/Collapse';
import {ResourceType, TDraftLinter} from 'types/Settings.types';
import {CollapsePanel} from 'components/Collapse';
import * as S from '../common/Settings.styled';
import PluginHeader from './PluginHeader';
import Plugin from './Plugin';

const FORM_ID = 'linter';

const LinterForm = () => {
  const [form] = Form.useForm<TDraftLinter>();
  const {isLoading, onSubmit} = useSettings();
  const {linter} = useSettingsValues();

  useEffect(() => {
    form.setFieldsValue(linter);
  }, [form, linter]);

  const handleOnSubmit = (values: TDraftLinter) => {
    onSubmit([
      SettingService.getDraftResource(ResourceType.AnalyzerType, {
        ...values,
        minimumScore: parseInt(String(values?.minimumScore ?? 0), 10),
      }),
    ]);
  };

  return (
    <Form<TDraftLinter>
      autoComplete="off"
      form={form}
      initialValues={linter}
      layout="vertical"
      name={FORM_ID}
      onFinish={handleOnSubmit}
    >
      <Form.Item hidden name="id" />
      <Form.Item hidden name="name" />

      <Row gutter={32} align="middle">
        <Col span={8}>
          <Form.Item
            label="Minimum score"
            name="minimumScore"
            rules={[{required: true, message: 'Minimum score is required'}]}
          >
            <Input suffix="%" placeholder="0 to 100" type="number" />
          </Form.Item>
        </Col>

        <Col span={8}>
          <S.SwitchContainer>
            <Form.Item name="enabled" valuePropName="checked" noStyle>
              <Switch />
            </Form.Item>
            <label htmlFor={`${FORM_ID}_enabled`}>Enable Analyzer for All Tests</label>
          </S.SwitchContainer>
        </Col>
      </Row>

      <Row>
        <Col span={16}>
          <Form.List name="plugins">
            {fields => (
              <Collapse>
                {fields.map(field => (
                  <CollapsePanel key={field.key} header={<PluginHeader fieldKey={field.name} />}>
                    <Plugin fieldKey={field.name} baseName={['plugins', `${field.name}`]} />
                  </CollapsePanel>
                ))}
              </Collapse>
            )}
          </Form.List>
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

export default LinterForm;
