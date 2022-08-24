import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import * as jsyaml from 'js-yaml';
import {OtelReference} from '../../components/AssertionForm/hooks/useGetOTELSemanticConvertionAttributesInfo';
import {CompleteAttribute, OTELYaml} from './OTELYaml';

const PATH = `https://raw.githubusercontent.com/open-telemetry/opentelemetry-specification/main/`;

enum Tags {
  TAGS = 'tags',
}

function normalizeThis(examples?: Array<string | number> | number | string): string[] {
  switch (typeof examples) {
    case 'object':
      return examples.map(d => d.toString());
    case 'number':
      return [examples.toString()];
    case 'string':
      return [examples];
    default:
      return [];
  }
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
          url: `semantic_conventions/${folder}/${file}.yaml`,
          responseHandler: 'text',
        };
      },
      providesTags: (result, error, {kind}) => [{type: Tags.TAGS, id: kind}],
      transformResponse: (rawSpanList: string) => {
        const message: OTELYaml = jsyaml.load(rawSpanList);
        return (
          (message?.groups || [])
            .flatMap<CompleteAttribute>(s => (s.attributes || []).map(d => ({...d, group: s.id})))
            .reduce((acc: OtelReference, d: CompleteAttribute) => {
              let id = `${d.group}.${d?.ref || d?.id || ''}`;
              acc[id] = {description: d?.brief || '', tags: normalizeThis(d?.examples)};
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
