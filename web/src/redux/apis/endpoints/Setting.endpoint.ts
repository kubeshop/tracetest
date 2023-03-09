import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import Config, {TRawConfig} from 'models/Config.model';
import Polling, {TRawPolling} from 'models/Polling.model';
import {EResourceType, TListResponse, TResource, TSpec} from 'types/Settings.types';
import {TTestApiEndpointBuilder} from 'types/Test.types';

const ConfigEndpoint = (builder: TTestApiEndpointBuilder) => ({
  getConfig: builder.query<Config, unknown>({
    query: () => ({
      url: '/config/current',
      method: HTTP_METHOD.GET,
      headers: {
        'content-type': 'application/json',
      },
    }),
    providesTags: () => [{type: TracetestApiTags.SETTING, id: EResourceType.Config}],
    transformResponse: (rawConfig: TResource<TRawConfig>) => Config(rawConfig),
  }),
  getPolling: builder.query<Polling, unknown>({
    query: () => ({
      url: '/polling/',
      method: HTTP_METHOD.GET,
      headers: {
        'content-type': 'application/json',
      },
    }),
    providesTags: () => [{type: TracetestApiTags.SETTING, id: EResourceType.Polling}],
    transformResponse: (rawPollings: TListResponse<TRawPolling>) => Polling(rawPollings),
  }),
  createSetting: builder.mutation<undefined, {resource: TResource<TSpec>}>({
    query: ({resource}) => ({
      url: `/${resource.type.toLowerCase()}/`,
      method: HTTP_METHOD.POST,
      body: resource,
    }),
    invalidatesTags: (result, error, args) => [{type: TracetestApiTags.SETTING, id: args.resource.type}],
  }),
  updateSetting: builder.mutation<undefined, {resource: TResource<TSpec>}>({
    query: ({resource}) => ({
      url: `/${resource.type.toLowerCase()}/${resource.spec.id}`,
      method: HTTP_METHOD.PUT,
      body: resource,
    }),
    invalidatesTags: (result, error, args) => [{type: TracetestApiTags.SETTING, id: args.resource.type}],
  }),
});

export default ConfigEndpoint;
