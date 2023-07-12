import Test from 'models/Test.model';
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
  format: CliCommandFormat;
  environmentId?: string;
  test: Test;
  fileName: string;
};

type TApplyProps = {
  command: string;
  test: Test;
  environmentId?: string;
  enabled: boolean;
  fileName: string;
};
type TApplyOption = (props: TApplyProps) => string;

const CliCommandService = () => ({
  applyOptions: {
    [CliCommandOption.UseId]: ({enabled, command, test: {id}, fileName}) =>
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
  getCommand({options, format, test, environmentId, fileName}: TCliCommandConfig) {
    const command = Object.entries(options).reduce(
      (acc, [option, enabled]) =>
        this.applyOptions[option as CliCommandOption]({command: acc, enabled, test, environmentId, fileName}),
      'run test'
    );

    return `${command} --output ${format}`;
  },
});

export default CliCommandService();
