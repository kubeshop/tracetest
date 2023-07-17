import {ResourceType} from 'types/Resource.type';
import {getServerBaseUrl} from 'utils/Common';

export enum CliCommandOption {
  UseId = 'useId',
  SkipResultWait = 'skipWait',
  UseHostname = 'useHostname',
  UseCurrentEnvironment = 'useCurrentEnvironment',
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
  environmentId?: string;
  fileName: string;
  format: CliCommandFormat;
  resourceType: ResourceType;
};

type TApplyProps = {
  command: string;
  id: string;
  environmentId?: string;
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
    [CliCommandOption.UseCurrentEnvironment]: ({command, enabled, environmentId}) =>
      enabled && environmentId ? `${command} --environment ${environmentId}` : command,
    [CliCommandOption.GeneratesJUnit]: ({command, enabled}) => (enabled ? `${command} --junit result.junit` : command),
    [CliCommandOption.useDocker]: ({enabled, command}) =>
      `${
        enabled
          ? 'docker run --rm -it -v$(pwd):$(pwd) -w $(pwd) --network host --entrypoint tracetest kubeshop/tracetest:latest -s http://localhost:11633/'
          : 'tracetest'
      } ${command}`,
  } as Record<CliCommandOption, TApplyOption>,
  getCommand({options, format, id, environmentId, fileName, resourceType}: TCliCommandConfig) {
    const command = Object.entries(options).reduce(
      (acc, [option, enabled]) =>
        this.applyOptions[option as CliCommandOption]({command: acc, enabled, id, environmentId, fileName}),
      `run ${resourceType}`
    );

    return `${command} --output ${format}`;
  },
});

export default CliCommandService();
