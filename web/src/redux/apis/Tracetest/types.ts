import {BaseQueryFn, FetchArgs, FetchBaseQueryError} from '@reduxjs/toolkit/dist/query';
import {EndpointBuilder} from '@reduxjs/toolkit/dist/query/endpointDefinitions';
import {TracetestApiTags} from '../../../constants/Test.constants';

export type TTestApiEndpointBuilder = EndpointBuilder<
  BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError>,
  TracetestApiTags,
  'tracetest'
>;
