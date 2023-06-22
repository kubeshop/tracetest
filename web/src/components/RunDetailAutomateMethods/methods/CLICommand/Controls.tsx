import {toUpper} from 'lodash';
import {useEffect} from 'react';
import {Form, Radio, Typography} from 'antd';
import {CliCommandFormat, CliCommandOption, TCliCommandConfig} from 'services/CliCommand.service';
import * as S from './CliCommand.styled';
import SwitchControl from './SwitchControl';
import {defaultOptions} from './hooks/useCliCommand';
import Test from '../../../../models/Test.model';

const optionToText: Record<CliCommandOption, string> = {
  [CliCommandOption.Wait]: 'Wait for test to complete',
  [CliCommandOption.UseHostname]: 'Specify Tracetest server hostname',
  [CliCommandOption.UseCurrentEnvironment]: 'Use current environment',
  // [CliCommandOption.UseId]: 'Use file or test id',
  [CliCommandOption.GeneratesJUnit]: 'Generate JUnit report',
  [CliCommandOption.useDocker]: 'Run CLI via Docker image',
};

interface IProps {
  onChange(cmdConfig: TCliCommandConfig): void;
  test: Test;
  environmentId?: string;
}

const Controls = ({onChange, test, environmentId}: IProps) => {
  const [form] = Form.useForm<TCliCommandConfig>();
  const options = Form.useWatch('options', form);
  const format = Form.useWatch('format', form);

  useEffect(() => {
    onChange({
      options: options ?? defaultOptions,
      format: format ?? CliCommandFormat.Pretty,
      test,
      environmentId,
    });
  }, [environmentId, format, onChange, options, test]);

  return (
    <Form<TCliCommandConfig>
      form={form}
      autoComplete="off"
      initialValues={{
        options: defaultOptions,
        format: CliCommandFormat.Pretty,
      }}
      layout="horizontal"
      name="DEEP_LINK"
    >
      <S.ControlsContainer>
        <S.OptionsContainer>
          {Object.entries(optionToText).map(([name, text]) => (
            <Form.Item name={['options', name]} noStyle>
              <SwitchControl id={name} text={text} key={name} />
            </Form.Item>
          ))}
        </S.OptionsContainer>
        <S.FormatContainer>
          <Typography.Paragraph>Output Format:</Typography.Paragraph>
          <Form.Item name="format" noStyle>
            <Radio.Group>
              {Object.values(CliCommandFormat).map(cmdFormat => (
                <Radio key={cmdFormat} value={cmdFormat}>
                  {toUpper(cmdFormat)}
                </Radio>
              ))}
            </Radio.Group>
          </Form.Item>
        </S.FormatContainer>
      </S.ControlsContainer>
    </Form>
  );
};

export default Controls;
