import {useMemo} from 'react';
import EditTestSuite from 'components/EditTestSuite';
import TestSuiteRunResult from 'components/TestSuiteRunResult';
import useDocumentTitle from 'hooks/useDocumentTitle';
import {useTestSuite} from 'providers/TestSuite/TestSuite.provider';
import {useTestSuiteRun} from 'providers/TestSuiteRun/TestSuite.provider';
import TestSuiteService from 'services/TestSuite.service';
import * as S from './TestSuiteRunOverview.styled';

const Content = () => {
  const {testSuite} = useTestSuite();
  const {run} = useTestSuiteRun();
  useDocumentTitle(`${testSuite.name} - ${run.state}`);
  const draft = useMemo(() => TestSuiteService.getInitialValues(testSuite), [testSuite]);

  return (
    <S.Container>
      <S.SectionLeft>
        <EditTestSuite testSuite={draft} testSuiteRun={run} />
      </S.SectionLeft>
      <S.SectionRight>
        <TestSuiteRunResult testSuite={testSuite} testSuiteRun={run} />
      </S.SectionRight>
    </S.Container>
  );
};

export default Content;
