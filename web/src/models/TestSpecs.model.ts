import {TRawTestSpecs, TTestSpecs} from 'types/TestSpecs.types';

const TestSpecs = ({specs = []}: TRawTestSpecs): TTestSpecs => ({
  specs: specs.map(({selector: {query = ''} = {}, assertions = []}) => ({
    assertions,
    isDeleted: false,
    isDraft: false,
    selector: query,
  })),
});

export default TestSpecs;
