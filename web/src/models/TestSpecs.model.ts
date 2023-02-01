import {Model, TTestSchemas} from 'types/Common.types';

export type TTestSpecEntry = {
  assertions: string[];
  isDeleted?: boolean;
  isDraft: boolean;
  originalSelector?: string;
  selector: string;
};

export type TRawTestSpecEntry = {
  selector: {query: string};
  assertions: string[];
};

export type TRawTestSpecs = TTestSchemas['TestSpecs'];
type TestSpecs = Model<TRawTestSpecs, {specs: TTestSpecEntry[]}>;

const TestSpecs = ({specs = []}: TRawTestSpecs): TestSpecs => ({
  specs: specs.map(({selector: {query = ''} = {}, assertions = []}) => ({
    assertions,
    isDeleted: false,
    isDraft: false,
    selector: query,
  })),
});

export default TestSpecs;
