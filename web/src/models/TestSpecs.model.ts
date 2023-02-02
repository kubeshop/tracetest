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

export type TRawTestSpecs = TTestSchemas['TestSpecs'];
type TestSpecs = Model<TRawTestSpecs, {specs: TTestSpecEntry[]}>;

const TestSpecs = ({specs = []}: TRawTestSpecs): TestSpecs => ({
  specs: specs.map(({selector: {query = ''} = {}, assertions = [], name = ''}) => ({
    assertions,
    isDeleted: false,
    isDraft: false,
    selector: query,
    name: name ?? '',
  })),
});

export default TestSpecs;
