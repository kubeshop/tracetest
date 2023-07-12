import {useCallback, useState} from 'react';
import CliCommandService, {
  CliCommandOption,
  TCliCommandConfig,
  TCliCommandEnabledOptions,
} from 'services/CliCommand.service';

export const defaultOptions: TCliCommandEnabledOptions = {
  [CliCommandOption.Wait]: false,
  [CliCommandOption.UseHostname]: false,
  [CliCommandOption.UseCurrentEnvironment]: true,
  // [CliCommandOption.UseId]: true,
  [CliCommandOption.GeneratesJUnit]: false,
  [CliCommandOption.useDocker]: false,
};

const useCliCommand = () => {
  const [command, setCommand] = useState<string>('');

  const onGetCommand = useCallback(
    (config: TCliCommandConfig) => {
      const cmd = CliCommandService.getCommand(config);
      setCommand(cmd);
    },
    []
  );

  return {command, onGetCommand};
};

export default useCliCommand;
