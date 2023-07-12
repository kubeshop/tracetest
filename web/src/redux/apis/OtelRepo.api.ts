import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import * as jsyaml from 'js-yaml';
import {OtelReference} from '../../components/TestSpecForm/hooks/useGetOTELSemanticConventionAttributesInfo';
import {CompleteAttribute, OTELYaml} from './OTELYaml';

const PATH = 'https://raw.githubusercontent.com/open-telemetry/semantic-conventions/main/';

enum Tags {
  TAGS = 'tags',
}

const OtelRepoAPI = createApi({
  reducerPath: 'otel',
  baseQuery: fetchBaseQuery({
    baseUrl: PATH,
  }),
  tagTypes: Object.values(Tags),
  endpoints: build => ({
    getConventions: build.query<OtelReference, {kind?: string; folder?: string}>({
      query: ({folder = 'trace', kind: file = 'http'}) => {
        return {
          url: `model/${folder}/${file}.yaml`,
          responseHandler: 'text',
        };
      },
      providesTags: (result, error, {kind}) => [{type: Tags.TAGS, id: kind}],
      transformResponse: (rawSpanList: string) => {
        const message: OTELYaml = jsyaml.load(rawSpanList);
        return (
          (message?.groups || [])
            .flatMap<CompleteAttribute>(s => (s.attributes || []).map(d => ({...d, group: s.prefix})))
            .reduce((acc: OtelReference, d: CompleteAttribute) => {
              let id = `${d.group}.${d?.ref || d?.id || ''}`;
              acc[id] = {description: d?.brief || '', note: d?.note ?? '', tags: []};
              return acc;
            }, {}) || {}
        );
      },
    }),
  }),
});

export const {useGetConventionsQuery} = OtelRepoAPI;
export const {endpoints} = OtelRepoAPI;

export default OtelRepoAPI;
