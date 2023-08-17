import {SortBy, SortDirection, TracetestApiTags} from 'constants/Test.constants';
import {PaginationResponse} from 'hooks/usePagination';
import Resource, {TRawResource} from 'models/Resource.model';
import {ResourceType} from 'types/Resource.type';
import {getTotalCountFromHeaders} from 'utils/Common';
import { TTestApiEndpointBuilder } from '../types';

const defaultHeaders = {'content-type': 'application/json', 'X-Tracetest-Augmented': 'true'};

export const resourceEndpoints = (builder: TTestApiEndpointBuilder) => ({
  getResources: builder.query<
    PaginationResponse<Resource>,
    {resourceType: ResourceType, take?: number; skip?: number; query?: string; sortBy?: SortBy; sortDirection?: SortDirection}
  >({
    query: ({resourceType, take = 25, skip = 0, query = '', sortBy = '', sortDirection = ''}) => ({
      url: `/${resourceType}?take=${take}&skip=${skip}&query=${query}&sortBy=${sortBy}&sortDirection=${sortDirection}`,
      headers: defaultHeaders,
    }),
    providesTags: () => [{type: TracetestApiTags.RESOURCE, id: 'LIST'}],
    transformResponse: (rawResources: TRawResource[], meta) => {
      return {
        items: rawResources.map(rawResource => Resource(rawResource)),
        total: getTotalCountFromHeaders(meta),
      };
    },
  }),
  getResourceDefinition: builder.query<string, {resourceId: string; version?: number; resourceType: ResourceType}>({
    query: ({resourceId, resourceType}) => ({
      url: `/${resourceType}/${resourceId}`,
      responseHandler: 'text',
      headers: {
        'content-type': 'text/yaml',
      },
    }),
    providesTags: (result, error, {resourceId, version}) => [
      {type: TracetestApiTags.RESOURCE, id: `${resourceId}-${version}-definition`},
    ],
  }),
});
