import {TEST_RUNNER_DOCUMENTATION_URL} from 'constants/Common.constants';
import DocsBanner from 'components/DocsBanner/DocsBanner';
import * as S from '../common/Settings.styled';
import TestRunnerForm from './TestRunnerForm';

const TestRunner = () => {
  return (
    <S.Container>
      <S.Description>
        The Test Runner allows you to configure the behavior used to execute your tests and generate the results
        <DocsBanner>
          Need more information about Test Runner?{' '}
          <a href={TEST_RUNNER_DOCUMENTATION_URL} target="_blank">
            Go to our docs
          </a>
        </DocsBanner>
      </S.Description>
      <S.FormContainer>
        <TestRunnerForm />
      </S.FormContainer>
    </S.Container>
  );
};

export default TestRunner;
