import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import {TTest} from 'types/Test.types';
import {useDeleteTestByIdMutation} from 'redux/apis/TraceTest.api';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {useConfirmationModal} from '../providers/ConfirmationModal/ConfirmationModal.provider';

const useDeleteTest = () => {
  const [deleteTestMutation] = useDeleteTestByIdMutation();
  const navigate = useNavigate();

  const {onOpen} = useConfirmationModal();

  const onConfirmDelete = useCallback(
    (test: TTest) => {
      TestAnalyticsService.onDeleteTest();
      deleteTestMutation({testId: test?.id || ''});
      navigate('/');
    },
    [deleteTestMutation, navigate]
  );

  const onDelete = useCallback(
    (test: TTest) => {
      onOpen(`Are you sure you want to delete “${test?.name}”?`, () => onConfirmDelete(test));
    },
    [onConfirmDelete, onOpen]
  );

  return onDelete;
};

export default useDeleteTest;
