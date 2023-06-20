import {useCallback, useState} from 'react';
import {useEnvironment} from 'providers/Environment/Environment.provider';
import {useTest} from 'providers/Test/Test.provider';
import CliCommandService, {TCliCommandConfig} from 'services/CliCommand.service';

const useCliCommand = () => {
  const {selectedEnvironment} = useEnvironment();
  const {test} = useTest();
  const [command, setCommand] = useState<string>('');

  const onGetCommand = useCallback(
    (config: TCliCommandConfig) => {
      const cmd = CliCommandService.getCommand(config, test, selectedEnvironment?.id);
      setCommand(cmd);
    },
    [selectedEnvironment?.id, test]
  );

  return {command, onGetCommand};
};

export default useCliCommand;
