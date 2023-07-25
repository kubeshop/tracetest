import {useCallback} from 'react';

import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import {useDeleteRunByIdMutation, useDeleteTransactionRunByIdMutation} from 'redux/apis/TraceTest.api';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {ResourceType} from 'types/Resource.type';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';

interface IProps {
  id: string;
  isRunView?: boolean;
  type: ResourceType;
}

const useDeleteResourceRun = ({id, isRunView = false, type}: IProps) => {
  const [deleteTestRunById] = useDeleteRunByIdMutation();
  const [deleteTransactionRunById] = useDeleteTransactionRunByIdMutation();
  const {navigate} = useDashboard();
  const {onOpen} = useConfirmationModal();

  const onConfirmDelete = useCallback(
    (runId: string) => {
      if (type === ResourceType.Test) {
        TestAnalyticsService.onDeleteTestRun();
        deleteTestRunById({testId: id, runId});
        if (isRunView) navigate(`/test/${id}`);
      } else if (type === ResourceType.Transaction) {
        deleteTransactionRunById({transactionId: id, runId});
        if (isRunView) navigate(`/transaction/${id}`);
      }
    },
    [deleteTestRunById, deleteTransactionRunById, id, isRunView, navigate, type]
  );

  return useCallback(
    (runId: string) => {
      onOpen({
        title: `Are you sure you want to delete the Run?`,
        onConfirm: () => onConfirmDelete(runId),
      });
    },
    [onConfirmDelete, onOpen]
  );
};

export default useDeleteResourceRun;
