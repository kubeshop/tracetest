import {TRawTestSpecs, TTestSpecs} from 'types/TestSpecs.types';
import Assertion from './Assertion.model';

const TestSpecs = ({specs = []}: TRawTestSpecs): TTestSpecs => ({
  specs: specs.map(({selector: {query = ''} = {}, assertions = []}) => ({
    assertions: assertions.map(rawAssertion => Assertion(rawAssertion)),
    isDeleted: false,
    isDraft: false,
    selector: query,
  })),
});

export default TestSpecs;
