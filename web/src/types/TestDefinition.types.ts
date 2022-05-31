import {TAssertion} from './Assertion.types';
import {Model, TTestSchemas} from './Common.types';

export type TRawTestDefinition = TTestSchemas['TestDefinition'];

export type TTestDefinitionEntry = {
  id?: string;
  selector: string;
  assertionList: TAssertion[];
  isDraft: boolean;
  isDeleted?: boolean;
};

export type TRawTestDefinitionEntry = {
  selector: string;
  assertions: TAssertion[];
};

export type TTestDefinition = Model<
  TRawTestDefinition,
  {
    definitionList: TTestDefinitionEntry[];
    definitions?: TRawTestDefinition;
  }
>;
