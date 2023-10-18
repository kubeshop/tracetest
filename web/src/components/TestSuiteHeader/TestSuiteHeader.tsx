import {useLocation} from 'react-router-dom';
import CreateButton from 'components/CreateButton';
import {TestSuiteRunStatusIcon} from 'components/RunStatusIcon';
import TestState from 'components/TestState';
import TestSuiteRunActionsMenu from 'components/TestSuiteRunActionsMenu';
import {TestState as TestStateEnum} from 'constants/TestRun.constants';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {useTestSuite} from 'providers/TestSuite/TestSuite.provider';
import {useTestSuiteRun} from 'providers/TestSuiteRun/TestSuite.provider';
import {useAppSelector} from 'redux/hooks';
import UserSelectors from 'selectors/User.selectors';
import * as S from './TestSuiteHeader.styled';
import VariableSetSelector from '../VariableSetSelector/VariableSetSelector';

const testSuiteLastPathRegex = /\/testsuite\/[\w-]+\/run\/[\d-]+\/([\w-]+)/;

function getLastPath(pathname: string): string {
  const match = pathname.match(testSuiteLastPathRegex);
  if (match === null) {
    return '';
  }

  return match.length > 1 ? match[1] : '';
}

const LINKS = [
  {id: 'overview', label: 'Overview'},
  {id: 'automate', label: 'Automate'},
];

const TestSuiteHeader = () => {
  const {testSuite, onRun} = useTestSuite();
  const {run} = useTestSuiteRun();
  const {navigate} = useDashboard();
  const {pathname} = useLocation();
  const {id: testSuiteId, name, version, description} = testSuite;
  const {state, id: runId, allStepsRequiredGatesPassed} = run;
  const lastPath = getLastPath(pathname);
  const runOriginPath = useAppSelector(UserSelectors.selectRunOriginPath);

  return (
    <S.Container>
      <S.Section>
        <a onClick={() => navigate(runOriginPath)} data-cy="testsuite-header-back-button">
          <S.BackIcon />
        </a>
        <div>
          <S.Title data-cy="testsuite-details-name">
            {name} (v{version})
          </S.Title>
          <S.Text>{description}</S.Text>
        </div>
      </S.Section>

      <S.LinksContainer>
        {LINKS.map(({id, label}) => (
          <S.Link
            key={id}
            to={`/testsuite/${testSuiteId}/run/${runId}/${id}`}
            $isActive={lastPath === id || (!lastPath && id === LINKS[0].id)}
          >
            {label}
          </S.Link>
        ))}
      </S.LinksContainer>

      <S.SectionRight>
        {state && state !== TestStateEnum.FINISHED && (
          <S.StateContainer data-cy="testsuite-run-result-status">
            <S.StateText>Status:</S.StateText>
            <TestState testState={state} />
          </S.StateContainer>
        )}
        {state && state === TestStateEnum.FINISHED && (
          <TestSuiteRunStatusIcon state={state!} hasFailedTests={!allStepsRequiredGatesPassed} />
        )}
        <VariableSetSelector />
        {state && state === TestStateEnum.FINISHED && (
          <CreateButton data-cy="testsuite-run-button" ghost onClick={() => onRun(runId)} type="primary">
            Run Test Suite
          </CreateButton>
        )}
        <TestSuiteRunActionsMenu testSuiteId={testSuiteId} runId={runId} isRunView />
      </S.SectionRight>
    </S.Container>
  );
};

export default TestSuiteHeader;
