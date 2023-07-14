import {Form, Radio, Select, Typography} from 'antd';
import {toUpper} from 'lodash';
import {useEffect, useMemo} from 'react';
import {SupportedRequiredGates} from 'models/TestRunner.model';
import {CliCommandFormat, CliCommandOption, TCliCommandConfig} from 'services/CliCommand.service';
import {ResourceType} from 'types/Resource.type';
import * as S from './CliCommand.styled';
import SwitchControl from './SwitchControl';
import {defaultOptions} from './hooks/useCliCommand';

interface IOptionsMetadataParams {
  isEnvironmentSelected: boolean;
  resourceType: ResourceType;
}
interface IOptionsMetadata {
  label: string;
  help?: string;
  disabled?: boolean;
}

function getOptionsMetadata({
  isEnvironmentSelected,
  resourceType,
}: IOptionsMetadataParams): Record<CliCommandOption, IOptionsMetadata> {
  return {
    [CliCommandOption.UseId]: {label: `Use ${resourceType} ID instead of file`},
    [CliCommandOption.SkipResultWait]: {label: `Skip waiting for ${resourceType} to complete`},
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
  id: string;
  environmentId?: string;
  fileName: string;
  resourceType: ResourceType;
}

const Controls = ({onChange, id, environmentId, fileName, resourceType}: IProps) => {
  const [form] = Form.useForm<TCliCommandConfig>();
  const options = Form.useWatch('options', form);
  const format = Form.useWatch('format', form);
  const requiredGates = Form.useWatch('required-gates', form);
  const optionsMetadata = useMemo(
    () => getOptionsMetadata({isEnvironmentSelected: !!environmentId, resourceType}),
    [environmentId, resourceType]
  );

  useEffect(() => {
    onChange({
      options: options ?? defaultOptions,
      format: format ?? CliCommandFormat.Pretty,
      requiredGates,
      id,
      environmentId,
      fileName,
      resourceType,
    });
  }, [environmentId, fileName, format, requiredGates, onChange, options, id, resourceType]);

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
            <Form.Item key={name} name={['options', name]} noStyle>
              <SwitchControl id={name} text={data.label} key={name} disabled={data.disabled} help={data.help} />
            </Form.Item>
          ))}
        </S.OptionsContainer>

        <S.RequiredGatesContainer>
          <Typography.Paragraph>Override default Required Gates:</Typography.Paragraph>
          <Form.Item name="required-gates">
            <Select mode="multiple" placeholder="Please select the required gates">
              {Object.values(SupportedRequiredGates).map(requiredGate => (
                <Select.Option key={requiredGate} value={requiredGate}>
                  {requiredGate}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
        </S.RequiredGatesContainer>

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
