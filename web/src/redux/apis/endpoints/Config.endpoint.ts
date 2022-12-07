import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import {TConfig, TRawConfig, TTestConnectionRequest, TTestConnectionResponse} from 'types/Config.types';
// import Config from 'models/Config.model';
import ConfigMock from 'models/__mocks__/Config.mock';

const ConfigEndpoint = (builder: TTestApiEndpointBuilder) => ({
  getConfig: builder.query<TConfig, unknown>({
    // query: () => '/config',
    query: () => '/tests',
    providesTags: () => [{type: TracetestApiTags.CONFIG, id: 'config'}],
    transformResponse: () => ConfigMock.model(),
  }),
  updateConfig: builder.mutation<undefined, TRawConfig>({
    query: config => ({
      url: '/config',
      method: HTTP_METHOD.PUT,
      body: config,
    }),
    invalidatesTags: [{type: TracetestApiTags.CONFIG, id: 'config'}],
  }),
  testConnection: builder.mutation<TTestConnectionResponse, TTestConnectionRequest>({
    query: connectionTest => ({
      url: `/config/connection`,
      method: HTTP_METHOD.PUT,
      body: connectionTest,
    }),
  }),
});

export default ConfigEndpoint;
