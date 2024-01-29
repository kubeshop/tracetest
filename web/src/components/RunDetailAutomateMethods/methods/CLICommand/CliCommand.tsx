import {withCustomization} from 'providers/Customization';
import useCliCommand from './hooks/useCliCommand';
import {IMethodChildrenProps} from '../../RunDetailAutomateMethods';
import Command from './Command';

const CLiCommand = (props: IMethodChildrenProps) => {
  const {command, onGetCommand} = useCliCommand();

  return <Command command={command} onGetCommand={onGetCommand} {...props} />;
};

export default withCustomization(CLiCommand, 'cliCommand');
