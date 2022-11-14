import {SortBy, SortDirection, TracetestApiTags} from 'constants/Test.constants';
import Resource from 'models/Resource.model';
import {PaginationResponse} from 'hooks/usePagination';
import {TRawResource, TResource} from 'types/Resource.type';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import {getTotalCountFromHeaders} from 'utils/Common';

const ResourceEndpoint = (builder: TTestApiEndpointBuilder) => ({
  getResources: builder.query<
    PaginationResponse<TResource>,
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
});

export default ResourceEndpoint;
