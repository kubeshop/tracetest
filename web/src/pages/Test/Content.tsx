import {Button} from 'antd';
import {useCallback, useMemo} from 'react';
import {useNavigate} from 'react-router-dom';

import {Steps} from 'components/GuidedTour/testDetailsStepList';
import PaginatedList from 'components/PaginatedList';
import TestRunCard from 'components/RunCard/TestRunCard';
import TestHeader from 'components/TestHeader';
import useDeleteResource from 'hooks/useDeleteResource';
import useTestCrud from 'providers/Test/hooks/useTestCrud';
import {useTest} from 'providers/Test/Test.provider';
import {useGetRunListQuery} from 'redux/apis/TraceTest.api';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {ResourceType} from 'types/Resource.type';
import {TTestRun} from 'types/TestRun.types';
import ExperimentalFeature from 'utils/ExperimentalFeature';
import * as S from './Test.styled';
import useDocumentTitle from '../../hooks/useDocumentTitle';

const Content = () => {
  const navigate = useNavigate();
  const {test} = useTest();
  const onDeleteResource = useDeleteResource();
  const {runTest, isLoadingRunTest} = useTestCrud();
  const params = useMemo(() => ({testId: test.id}), [test.id]);
  useDocumentTitle(`${test.name}`);

  const handleRunTest = useCallback(async () => {
    if (test.id) runTest(test.id);
  }, [runTest, test.id]);

  return (
    <S.Container $isWhite={!ExperimentalFeature.isEnabled('transactions')}>
      <TestHeader
        description={`${test.trigger.type.toUpperCase()} • ${test.trigger.method.toUpperCase()} • ${
          test.trigger.entryPoint
        }`}
        id={test.id}
        onBack={() => navigate('/')}
        onDelete={() => onDeleteResource(test.id, test.name, ResourceType.Test)}
        title={`${test.name} (v${test.version})`}
      />

      <S.ActionsContainer>
        <div />
        <Button
          data-cy="test-details-run-test-button"
          data-tour={GuidedTourService.getStep(GuidedTours.TestDetails, Steps.RunTest)}
          ghost
          loading={isLoadingRunTest}
          onClick={handleRunTest}
          type="primary"
        >
          Run Test
        </Button>
      </S.ActionsContainer>

      <PaginatedList<TTestRun, {testId: string}>
        dataCy="run-card-list"
        itemComponent={({item}) => (
          <TestRunCard linkTo={`/test/${test.id}/run/${item.id}`} run={item} testId={test.id} />
        )}
        params={params}
        query={useGetRunListQuery}
      />
    </S.Container>
  );
};

export default Content;
