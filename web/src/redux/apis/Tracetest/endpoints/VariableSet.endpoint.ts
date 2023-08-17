import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {PaginationResponse} from 'hooks/usePagination';
import VariableSet, {TRawVariableSet} from 'models/VariableSet.model';
import {TListResponse} from 'types/Settings.types';
import {TTestApiEndpointBuilder} from '../types';

export const variableSetEndpoints = (builder: TTestApiEndpointBuilder) => ({
  getVariableSets: builder.query<PaginationResponse<VariableSet>, {take?: number; skip?: number; query?: string}>({
    query: ({take = 25, skip = 0, query = ''}) => ({
      url: `/variablesets?take=${take}&skip=${skip}&query=${query}`,
      headers: {
        'content-type': 'application/json',
      },
    }),
    providesTags: () => [{type: TracetestApiTags.VARIABLE_SET, id: 'LIST'}],
    transformResponse: ({items, count}: TListResponse<TRawVariableSet>) => ({
      items: items.map(rawEnv => VariableSet(rawEnv)),
      total: count,
    }),
  }),
  createVariableSet: builder.mutation<undefined, TRawVariableSet>({
    query: variableSet => ({
      url: '/variablesets',
      method: HTTP_METHOD.POST,
      body: variableSet,
      headers: {
        'content-type': 'application/json',
      },
    }),
    invalidatesTags: [{type: TracetestApiTags.VARIABLE_SET, id: 'LIST'}],
  }),
  updateVariableSet: builder.mutation<undefined, {variableSet: TRawVariableSet; variableSetId: string}>({
    query: ({variableSet, variableSetId}) => ({
      url: `/variablesets/${variableSetId}`,
      method: HTTP_METHOD.PUT,
      body: variableSet,
      headers: {
        'content-type': 'application/json',
      },
    }),
    invalidatesTags: [{type: TracetestApiTags.VARIABLE_SET, id: 'LIST'}],
  }),
  deleteVariableSet: builder.mutation<undefined, {variableSetId: string}>({
    query: ({variableSetId}) => ({
      url: `/variablesets/${variableSetId}`,
      method: HTTP_METHOD.DELETE,
      headers: {
        'content-type': 'application/json',
      },
    }),
    invalidatesTags: [{type: TracetestApiTags.VARIABLE_SET, id: 'LIST'}],
  }),
});
