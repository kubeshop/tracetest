import {useCallback, useEffect, useState} from 'react';
import {toUpper} from 'lodash';
import {Radio, Typography} from 'antd';
import {
  CliCommandFormat,
  CliCommandOption,
  TCliCommandConfig,
  TCliCommandEnabledOptions,
} from 'services/CliCommand.service';
import * as S from './CliCommand.styled';
import SwitchControl from './SwitchControl';

const optionToText: Record<CliCommandOption, string> = {
  [CliCommandOption.Wait]: 'Wait for test to complete',
  [CliCommandOption.UseHostname]: 'Specify Tracetest server hostname',
  [CliCommandOption.UseCurrentEnvironment]: 'Use current environment',
  // [CliCommandOption.UseId]: 'Use file or test id',
  [CliCommandOption.GeneratesJUnit]: 'Generate JUnit report',
  [CliCommandOption.useDocker]: 'Run CLI via Docker image',
};

const defaultOptions: TCliCommandEnabledOptions = {
  [CliCommandOption.Wait]: false,
  [CliCommandOption.UseHostname]: false,
  [CliCommandOption.UseCurrentEnvironment]: true,
  // [CliCommandOption.UseId]: true,
  [CliCommandOption.GeneratesJUnit]: false,
  [CliCommandOption.useDocker]: false,
};

interface IProps {
  onChange(cmdConfig: TCliCommandConfig): void;
}

const Controls = ({onChange}: IProps) => {
  const [enabledOptions, setEnabledOptions] = useState<Record<CliCommandOption, boolean>>(defaultOptions);
  const [format, setFormat] = useState<CliCommandFormat>(CliCommandFormat.Pretty);

  const onSwitchChange = useCallback((name: CliCommandOption, isEnabled: boolean) => {
    setEnabledOptions(prev => ({...prev, [name]: isEnabled}));
  }, []);

  useEffect(() => {
    onChange({options: enabledOptions, format});
  }, [format, enabledOptions, onChange]);

  return (
    <S.ControlsContainer>
      <S.OptionsContainer>
        {Object.entries(optionToText).map(([name, text]) => (
          <SwitchControl
            text={text}
            key={name}
            onChange={isEnabled => onSwitchChange(name as CliCommandOption, isEnabled)}
            value={enabledOptions[name as CliCommandOption]}
          />
        ))}
      </S.OptionsContainer>
      <S.FormatContainer>
        <Typography.Paragraph>Output Format:</Typography.Paragraph>
        <Radio.Group defaultValue={format} value={format}>
          {Object.values(CliCommandFormat).map(cmdFormat => (
            <Radio key={cmdFormat} value={cmdFormat} onChange={() => setFormat(cmdFormat)}>
              {toUpper(cmdFormat)}
            </Radio>
          ))}
        </Radio.Group>
      </S.FormatContainer>
    </S.ControlsContainer>
  );
};

export default Controls;
