import {Form, Radio, Typography} from 'antd';
import {toUpper} from 'lodash';
import {useEffect, useMemo} from 'react';
import Test from 'models/Test.model';
import {CliCommandFormat, CliCommandOption, TCliCommandConfig} from 'services/CliCommand.service';
import * as S from './CliCommand.styled';
import SwitchControl from './SwitchControl';
import {defaultOptions} from './hooks/useCliCommand';

interface IOptionsMetadataParams {
  isEnvironmentSelected: boolean;
}
interface IOptionsMetadata {
  label: string;
  help?: string;
  disabled?: boolean;
}

function getOptionsMetadata({
  isEnvironmentSelected,
}: IOptionsMetadataParams): Record<CliCommandOption, IOptionsMetadata> {
  return {
    [CliCommandOption.UseId]: {label: 'Use test ID instead of file'},
    [CliCommandOption.SkipResultWait]: {label: 'Skip waiting for test to complete'},
    [CliCommandOption.UseHostname]: {label: 'Specify Tracetest server hostname'},
    [CliCommandOption.UseCurrentEnvironment]: {
      label: 'Use selected environment',
      help: !isEnvironmentSelected ? 'This option is only available when an environment is selected' : undefined,
      disabled: !isEnvironmentSelected,
    },
    [CliCommandOption.GeneratesJUnit]: {label: 'Generate JUnit report'},
    [CliCommandOption.useDocker]: {label: 'Run CLI via Docker image'},
  };
}

interface IProps {
  onChange(cmdConfig: TCliCommandConfig): void;
  test: Test;
  environmentId?: string;
  fileName: string;
}

const Controls = ({onChange, test, environmentId, fileName}: IProps) => {
  const [form] = Form.useForm<TCliCommandConfig>();
  const options = Form.useWatch('options', form);
  const format = Form.useWatch('format', form);
  const optionsMetadata = useMemo(() => getOptionsMetadata({isEnvironmentSelected: !!environmentId}), [environmentId]);

  useEffect(() => {
    onChange({
      options: options ?? defaultOptions,
      format: format ?? CliCommandFormat.Pretty,
      test,
      environmentId,
      fileName,
    });
  }, [environmentId, fileName, format, onChange, options, test]);

  return (
    <Form<TCliCommandConfig>
      form={form}
      autoComplete="off"
      initialValues={{
        options: defaultOptions,
        format: CliCommandFormat.Pretty,
      }}
      layout="horizontal"
      name="CLI_COMMAND"
    >
      <S.ControlsContainer>
        <S.OptionsContainer>
          {Object.entries(optionsMetadata).map(([name, data]) => (
            <Form.Item name={['options', name]} noStyle>
              <SwitchControl id={name} text={data.label} key={name} disabled={data.disabled} help={data.help} />
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
