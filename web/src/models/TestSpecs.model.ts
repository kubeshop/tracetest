import {Model, TTestSchemas} from 'types/Common.types';

export type TTestSpecEntry = {
  assertions: string[];
  isDeleted?: boolean;
  isDraft: boolean;
  originalSelector?: string;
  selector: string;
  name: string;
};

export type TRawTestSpecEntry = {
  selector: {query: string};
  assertions: string[];
  name: string;
};

// TechDebt: this is a temporary method to deal with the changes on the OpenAPI
// later we need to think and update our type system to deal with that
export const rawTestSpecToNewFormat = (spec : TRawTestSpecEntry) => {
  return {
    name: spec.name,
    assertions: spec.assertions,
    selector: spec.selector?.query,
    selectorParsed: spec.selector,
  };
};

export type TRawTestSpecs = TTestSchemas['TestSpecs'];
type TestSpecs = Model<TRawTestSpecs, {specs: TTestSpecEntry[]}>;

const TestSpecs = ({specs = []}: TRawTestSpecs): TestSpecs => {
  const newSpecs = specs.map(({selectorParsed: {query = ''} = {}, assertions = [], name = ''}) => ({
    assertions,
    isDeleted: false,
    isDraft: false,
    selector: query,
    name: name ?? '',
  }));

  return {
    specs: newSpecs
  };
};

export default TestSpecs;
