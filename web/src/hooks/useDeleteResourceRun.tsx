import {useCallback} from 'react';

import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import TracetestAPI from 'redux/apis/Tracetest';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {ResourceType} from 'types/Resource.type';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';

const {useDeleteRunByIdMutation, useDeleteTestSuiteRunByIdMutation} = TracetestAPI.instance;

interface IProps {
  id: string;
  isRunView?: boolean;
  type: ResourceType;
}

const useDeleteResourceRun = ({id, isRunView = false, type}: IProps) => {
  const [deleteTestRunById] = useDeleteRunByIdMutation();
  const [deleteTestSuiteRunById] = useDeleteTestSuiteRunByIdMutation();
  const {navigate} = useDashboard();
  const {onOpen} = useConfirmationModal();

  const onConfirmDelete = useCallback(
    (runId: number) => {
      if (type === ResourceType.Test) {
        TestAnalyticsService.onDeleteTestRun();
        deleteTestRunById({testId: id, runId});
        if (isRunView) navigate(`/test/${id}`);
      } else if (type === ResourceType.TestSuite) {
        deleteTestSuiteRunById({testSuiteId: id, runId});
        if (isRunView) navigate(`/testsuite/${id}`);
      }
    },
    [deleteTestRunById, deleteTestSuiteRunById, id, isRunView, navigate, type]
  );

  return useCallback(
    (runId: number) => {
      onOpen({
        title: `Are you sure you want to delete the Run?`,
        onConfirm: () => onConfirmDelete(runId),
      });
    },
    [onConfirmDelete, onOpen]
  );
};

export default useDeleteResourceRun;
