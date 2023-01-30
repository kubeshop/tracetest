import {TRawTestSpecs, TTestSpecs} from 'types/TestSpecs.types';

const TestSpecs = ({specs = []}: TRawTestSpecs): TTestSpecs => ({
  specs: specs.map(({selector: {query = ''} = {}, assertions = [], name = ''}) => ({
    assertions,
    isDeleted: false,
    isDraft: false,
    selector: query,
    name: name ?? '',
  })),
});

export default TestSpecs;
