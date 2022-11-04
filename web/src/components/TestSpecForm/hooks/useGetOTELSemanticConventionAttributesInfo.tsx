import AttributesTags from 'constants/AttributesTags.json';
import {useGetConventionsQuery} from 'redux/apis/OtelRepo.api';

export type OtelReference = Record<string, OtelReferenceModel>;

export interface OtelReferenceModel {
  description: string;
  tags: string[];
}

const attributesTags: OtelReference = AttributesTags;

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
    ...attributesTags,
  };
};
