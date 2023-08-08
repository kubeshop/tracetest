import {Form, Radio, Typography} from 'antd';
import {toUpper} from 'lodash';
import {useEffect, useMemo} from 'react';
import RequiredGatesInput from 'components/Settings/TestRunner/RequiredGatesInput';
import {TooltipQuestion} from 'components/TooltipQuestion/TooltipQuestion';
import {CliCommandFormat, CliCommandOption, TCliCommandConfig} from 'services/CliCommand.service';
import {ResourceType} from 'types/Resource.type';
import * as S from './CliCommand.styled';
import SwitchControl from './SwitchControl';
import {defaultOptions} from './hooks/useCliCommand';

interface IOptionsMetadataParams {
  isVariableSetSelected: boolean;
  resourceType: ResourceType;
}
interface IOptionsMetadata {
  label: string;
  help?: string;
  disabled?: boolean;
}

function getOptionsMetadata({
  isVariableSetSelected,
  resourceType,
}: IOptionsMetadataParams): Record<CliCommandOption, IOptionsMetadata> {
  return {
    [CliCommandOption.UseId]: {label: `Use ${resourceType} ID instead of file`},
    [CliCommandOption.SkipResultWait]: {label: `Skip waiting for ${resourceType} to complete`},
    [CliCommandOption.UseHostname]: {label: 'Specify Tracetest server hostname'},
    [CliCommandOption.UseCurrentVariableSet]: {
      label: 'Use selected variable set',
      help: !isVariableSetSelected ? 'This option is only available when a variable set is selected' : undefined,
      disabled: !isVariableSetSelected,
    },
    [CliCommandOption.GeneratesJUnit]: {label: 'Generate JUnit report'},
    [CliCommandOption.useDocker]: {label: 'Run CLI via Docker image'},
  };
}

interface IProps {
  onChange(cmdConfig: TCliCommandConfig): void;
  id: string;
  variableSetId?: string;
  fileName: string;
  resourceType: ResourceType;
}

const Controls = ({onChange, id, variableSetId, fileName, resourceType}: IProps) => {
  const [form] = Form.useForm<TCliCommandConfig>();
  const options = Form.useWatch('options', form);
  const format = Form.useWatch('format', form);
  const requiredGates = Form.useWatch('required-gates', form);
  const optionsMetadata = useMemo(
    () => getOptionsMetadata({isVariableSetSelected: !!variableSetId, resourceType}),
    [variableSetId, resourceType]
  );

  useEffect(() => {
    onChange({
      options: options ?? defaultOptions,
      format: format ?? CliCommandFormat.Pretty,
      requiredGates,
      id,
      variableSetId,
      fileName,
      resourceType,
    });
  }, [fileName, format, requiredGates, onChange, options, id, resourceType, variableSetId]);

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

        <Form.Item name="required-gates">
          <RequiredGatesInput
            title={
              <Typography.Paragraph>
                Override default Required Gates:{' '}
                <TooltipQuestion
                  margin={6}
                  title="Required Gates are used by the test runner to evaluate if a test is failed or not. You can override the default Required Gates for this run"
                />
              </Typography.Paragraph>
            }
          />
        </Form.Item>

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
