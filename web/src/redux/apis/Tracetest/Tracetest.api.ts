import {BaseQueryFn, FetchArgs, FetchBaseQueryError, createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import {TracetestApiTags} from 'constants/Test.constants';
import Env from 'utils/Env';

const customBaseQuery: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError> = async (
  args,
  api,
  extraOptions
) => {
  const baseUrl = Env.get('baseApiUrl');

  return fetchBaseQuery({
    baseUrl,
    prepareHeaders: headers => {
      headers.set('x-source', 'web');
    },
  })(args, api, extraOptions);
};

const TracetestAPI = createApi({
  reducerPath: 'tracetest',
  baseQuery: customBaseQuery,
  tagTypes: Object.values(TracetestApiTags),
  endpoints: () => ({}),
});

export default TracetestAPI;
