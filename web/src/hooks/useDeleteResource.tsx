import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import {useDeleteTestByIdMutation, useDeleteTransactionByIdMutation} from 'redux/apis/TraceTest.api';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import {ResourceType} from 'types/Resource.type';

const useDeleteResource = () => {
  const [deleteTestMutation] = useDeleteTestByIdMutation();
  const [deleteTransactionMutation] = useDeleteTransactionByIdMutation();
  const navigate = useNavigate();

  const {onOpen} = useConfirmationModal();

  const onConfirmDelete = useCallback(
    (id: string, type: ResourceType) => {
      if (type === ResourceType.Test) {
        TestAnalyticsService.onDeleteTest();
        deleteTestMutation({testId: id});
      } else if (type === ResourceType.Transaction) {
        deleteTransactionMutation({transactionId: id});
      }
      navigate('/');
    },
    [deleteTestMutation, deleteTransactionMutation, navigate]
  );

  const onDelete = useCallback(
    (id: string, name: string, type: ResourceType) => {
      onOpen(`Are you sure you want to delete “${name}”?`, () => onConfirmDelete(id, type));
    },
    [onConfirmDelete, onOpen]
  );

  return onDelete;
};

export default useDeleteResource;
