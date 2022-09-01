import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import {useDeleteRunByIdMutation} from 'redux/apis/TraceTest.api';

interface IProps {
  isRunView?: boolean;
  testId: string;
}

const useDeleteTestRun = ({isRunView = false, testId}: IProps) => {
  const [deleteRunById] = useDeleteRunByIdMutation();
  const navigate = useNavigate();
  const {onOpen} = useConfirmationModal();

  const onConfirmDelete = useCallback(
    (runId: string) => {
      TestAnalyticsService.onDeleteTestRun();
      deleteRunById({testId, runId});
      if (isRunView) navigate(`/test/${testId}`);
    },
    [deleteRunById, isRunView, navigate, testId]
  );

  const onDelete = useCallback(
    (runId: string) => {
      onOpen(`Are you sure you want to delete the Test Run?`, () => onConfirmDelete(runId));
    },
    [onConfirmDelete, onOpen]
  );

  return onDelete;
};

export default useDeleteTestRun;
