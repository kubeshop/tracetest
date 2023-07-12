import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import Config, {TRawConfig, TRawLiveConfig} from 'models/Config.model';
import Demo, {TRawDemo} from 'models/Demo.model';
import Linter, {TRawLinter} from 'models/Linter.model';
import Polling, {TRawPolling} from 'models/Polling.model';
import {ResourceType, TDraftResource, TListResponse} from 'types/Settings.types';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import {IListenerFunction} from 'gateways/WebSocket.gateway';
import WebSocketService from 'services/WebSocket.service';

const ConfigEndpoint = (builder: TTestApiEndpointBuilder) => ({
  getConfig: builder.query<Config, unknown>({
    query: () => ({
      url: '/configs/current',
      method: HTTP_METHOD.GET,
      headers: {
        'content-type': 'application/json',
      },
    }),
    providesTags: () => [{type: TracetestApiTags.SETTING, id: ResourceType.ConfigType}],
    transformResponse: (rawConfig: TRawConfig) => Config(rawConfig),
    async onCacheEntryAdded(arg, {cacheDataLoaded, cacheEntryRemoved, updateCachedData}) {
      const listener: IListenerFunction<TRawLiveConfig> = data => {
        updateCachedData(() => Config.FromLiveUpdate(data.event));
      };
      await WebSocketService.initWebSocketSubscription({
        listener,
        resource: '/app/config/update',
        waitToCleanSubscription: cacheEntryRemoved,
        waitToInitSubscription: cacheDataLoaded,
      });
    },
  }),
  getPolling: builder.query<Polling, unknown>({
    query: () => ({
      url: '/pollingprofiles/current',
      method: HTTP_METHOD.GET,
      headers: {
        'content-type': 'application/json',
      },
    }),
    providesTags: () => [{type: TracetestApiTags.SETTING, id: ResourceType.PollingProfileType}],
    transformResponse: (rawPolling: TRawPolling) => Polling(rawPolling),
  }),
  getDemo: builder.query<Demo[], unknown>({
    query: () => ({
      url: '/demos',
      method: HTTP_METHOD.GET,
      headers: {
        'content-type': 'application/json',
      },
    }),
    providesTags: () => [{type: TracetestApiTags.SETTING, id: ResourceType.DemoType}],
    transformResponse: ({items = []}: TListResponse<TRawDemo>) => items.map(rawDemo => Demo(rawDemo)),
  }),
  getLinter: builder.query<Linter, unknown>({
    query: () => ({
      url: '/analyzers/current',
      method: HTTP_METHOD.GET,
      headers: {'content-type': 'application/json'},
    }),
    providesTags: () => [{type: TracetestApiTags.SETTING, id: ResourceType.AnalyzerType}],
    transformResponse: (rawLinter: TRawLinter) => Linter(rawLinter),
  }),
  createSetting: builder.mutation<undefined, {resource: TDraftResource}>({
    query: ({resource}) => ({
      url: `/${resource.typePlural?.toLowerCase()}`,
      method: HTTP_METHOD.POST,
      body: resource,
    }),
    invalidatesTags: (result, error, args) => [{type: TracetestApiTags.SETTING, id: args.resource.type}],
  }),
  updateSetting: builder.mutation<undefined, {resource: TDraftResource}>({
    query: ({resource}) => ({
      url: `/${resource.typePlural?.toLowerCase()}/${resource.spec.id}`,
      method: HTTP_METHOD.PUT,
      body: resource,
    }),
    invalidatesTags: (result, error, args) => [{type: TracetestApiTags.SETTING, id: args.resource.type}],
  }),
});

export default ConfigEndpoint;
