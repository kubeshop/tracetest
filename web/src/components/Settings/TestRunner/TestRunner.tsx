import {TEST_RUNNER_DOCUMENTATION_URL} from 'constants/Common.constants';
import DocsBanner from 'components/DocsBanner/DocsBanner';
import * as S from '../common/Settings.styled';
import TestRunnerForm from './TestRunnerForm';

const TestRunner = () => {
  return (
    <S.Container>
      <S.Title level={2}>Test Runner</S.Title>
      <S.Description>
        <p>The Test Runner allows you to configure the behavior used to execute your tests and generate the results.</p>
        <DocsBanner>
          Need more information about the Test Runner?{' '}
          <a href={TEST_RUNNER_DOCUMENTATION_URL} target="_blank">
            Learn more in our docs
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
