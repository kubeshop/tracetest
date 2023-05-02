import {SortBy, SortDirection, TracetestApiTags} from 'constants/Test.constants';
import Resource, {TRawResource} from 'models/Resource.model';
import {PaginationResponse} from 'hooks/usePagination';
import {ResourceType} from 'types/Resource.type';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import {getTotalCountFromHeaders} from 'utils/Common';

const ResourceEndpoint = (builder: TTestApiEndpointBuilder) => ({
  getResources: builder.query<
    PaginationResponse<Resource>,
    {take?: number; skip?: number; query?: string; sortBy?: SortBy; sortDirection?: SortDirection}
  >({
    query: ({take = 25, skip = 0, query = '', sortBy = '', sortDirection = ''}) =>
      `/resources?take=${take}&skip=${skip}&query=${query}&sortBy=${sortBy}&sortDirection=${sortDirection}`,
    providesTags: () => [{type: TracetestApiTags.RESOURCE, id: 'LIST'}],
    transformResponse: (rawResources: TRawResource[], meta) => {
      return {
        items: rawResources.map(rawResource => Resource(rawResource)),
        total: getTotalCountFromHeaders(meta),
      };
    },
  }),
  getResourceDefinition: builder.query<string, {resourceId: string; version?: number; resourceType: ResourceType}>({
    query: ({resourceId, resourceType, version}) => ({
      url: `/${resourceType}s/${resourceId}${version ? `/version/${version}` : ''}/definition.yaml`,
      responseHandler: 'text',
    }),
    providesTags: (result, error, {resourceId, version}) => [
      {type: TracetestApiTags.RESOURCE, id: `${resourceId}-${version}-definition`},
    ],
  }),
  getResourceDefinitionV2: builder.query<string, {resourceId: string; version?: number; resourceType: ResourceType}>({
    query: ({resourceId, resourceType}) => ({
      url: `/${resourceType}s/${resourceId}`,
      responseHandler: 'text',
      headers: {
        'content-type': 'text/yaml',
      }
    }),
    providesTags: (result, error, {resourceId, version}) => [
      {type: TracetestApiTags.RESOURCE, id: `${resourceId}-${version}-definition`},
    ],
  }),
});

export default ResourceEndpoint;
