// import TracetestSemanticConventions from 'constants/TracetestSemanticConventions.json';
import {useGetConventionsQuery} from 'redux/apis/OtelRepo.api';

export type OtelReference = Record<string, OtelReferenceModel>;

export interface OtelReferenceModel {
  description: string;
  note: string;
  tags: string[];
}

// const tracetestSemanticConventions: OtelReference = TracetestSemanticConventions;

export const useGetOTELSemanticConventionAttributesInfo = (): OtelReference => {
  return {
    ...useGetConventionsQuery({kind: 'http'})?.data,
    ...useGetConventionsQuery({kind: 'database'})?.data,
    ...useGetConventionsQuery({kind: 'cloudevents'})?.data,
    ...useGetConventionsQuery({kind: 'compatibility'})?.data,
    ...useGetConventionsQuery({kind: 'trace-exception'})?.data,
    ...useGetConventionsQuery({kind: 'faas'})?.data,
    ...useGetConventionsQuery({kind: 'general'})?.data,
    ...useGetConventionsQuery({kind: 'messaging'})?.data,
    ...useGetConventionsQuery({kind: 'rpc'})?.data,
    // ...tracetestSemanticConventions,
  };
};
