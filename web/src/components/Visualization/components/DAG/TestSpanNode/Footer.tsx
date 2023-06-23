import AssertionResultChecks from 'components/AssertionResultChecks';
import TestOutputMark from 'components/TestOutputMark';
import TestOutput from 'models/TestOutput.model';
import TestRunOutput from 'models/TestRunOutput.model';
import {TTestSpecSummary} from 'types/TestRun.types';

interface IProps {
  testOutputs?: TestRunOutput[];
  testSpecs?: TTestSpecSummary;
}

const Footer = ({testOutputs, testSpecs}: IProps) => (
  <>
    {!!testOutputs?.length && <TestOutputMark outputs={testOutputs as TestOutput[]} />}
    {!!testSpecs && <AssertionResultChecks failed={testSpecs.failed} passed={testSpecs.passed} styleType="node" />}
  </>
);

export default Footer;
