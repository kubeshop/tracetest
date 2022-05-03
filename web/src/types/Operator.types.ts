import { Schemas } from "./Common.types";

export type TCompareOperatorName = 'eq' | 'ne' | 'gt' | 'lt' | 'gte' | 'lte' | 'contains';

export type TCompareOperatorSymbol = '==' | '<' | '>' | '!=' | '>=' | '<=' | 'contains';
export type TCompareOperator = Schemas['SpanAssertion']['operator'];