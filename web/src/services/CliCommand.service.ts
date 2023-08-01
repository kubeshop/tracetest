import {ResourceType} from 'types/Resource.type';
import {getServerBaseUrl} from 'utils/Common';

export enum CliCommandOption {
  UseId = 'useId',
  SkipResultWait = 'skipWait',
  UseHostname = 'useHostname',
  UseCurrentVariableSet = 'useCurrentVariableSet',
  GeneratesJUnit = 'generateJUnit',
  useDocker = 'useDocker',
}

export enum CliCommandFormat {
  Pretty = 'pretty',
  Json = 'json',
}

export type TCliCommandEnabledOptions = Record<CliCommandOption, boolean>;

export type TCliCommandConfig = {
  options: TCliCommandEnabledOptions;
  id: string;
  variableSetId?: string;
  fileName: string;
  format: CliCommandFormat;
  requiredGates: string[];
  resourceType: ResourceType;
};

type TApplyProps = {
  command: string;
  id: string;
  variableSetId?: string;
  enabled: boolean;
  fileName: string;
};
type TApplyOption = (props: TApplyProps) => string;

const CliCommandService = () => ({
  applyOptions: {
    [CliCommandOption.UseId]: ({enabled, command, id, fileName}) =>
      `${command} ${enabled ? `--id ${id}` : `--file ${fileName}`}`,
    [CliCommandOption.SkipResultWait]: ({command, enabled}) => (enabled ? `${command} --skip-result-wait` : command),
    [CliCommandOption.UseHostname]: ({command, enabled}) => {
      const baseUrl = getServerBaseUrl();
      return enabled ? `${command} --server-url ${baseUrl}` : command;
    },
    [CliCommandOption.UseCurrentVariableSet]: ({command, enabled, variableSetId}) =>
      enabled && variableSetId ? `${command} --variable_set ${variableSetId}` : command,
    [CliCommandOption.GeneratesJUnit]: ({command, enabled}) => (enabled ? `${command} --junit result.junit` : command),
    [CliCommandOption.useDocker]: ({enabled, command}) =>
      `${
        enabled
          ? 'docker run --rm -it -v$(pwd):$(pwd) -w $(pwd) --network host --entrypoint tracetest kubeshop/tracetest:latest -s http://localhost:11633/'
          : 'tracetest'
      } ${command}`,
  } as Record<CliCommandOption, TApplyOption>,

  applyRequiredGates: (command: string, requiredGates: string[]) =>
    requiredGates?.length ? `${command} --required-gates ${requiredGates.join(',')}` : command,

  getCommand({options, format, id, variableSetId, fileName, requiredGates, resourceType}: TCliCommandConfig) {
    let command = Object.entries(options).reduce(
      (acc, [option, enabled]) =>
        this.applyOptions[option as CliCommandOption]({command: acc, enabled, id, variableSetId, fileName}),
      `run ${resourceType}`
    );

    command = this.applyRequiredGates(command, requiredGates);

    return `${command} --output ${format}`;
  },
});

export default CliCommandService();
