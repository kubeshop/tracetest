import {useCallback} from 'react';
import {capitalize} from 'lodash';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import {useDeleteTestByIdMutation, useDeleteTransactionByIdMutation} from 'redux/apis/Tracetest';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {ResourceType} from 'types/Resource.type';
import {useNotification} from 'providers/Notification/Notification.provider';

const useDeleteResource = () => {
  const [deleteTestMutation] = useDeleteTestByIdMutation();
  const [deleteTransactionMutation] = useDeleteTransactionByIdMutation();
  const {navigate} = useDashboard();
  const {showNotification} = useNotification();

  const {onOpen} = useConfirmationModal();

  const onConfirmDelete = useCallback(
    async (id: string, type: ResourceType) => {
      try {
        if (type === ResourceType.Test) {
          TestAnalyticsService.onDeleteTest();
          await deleteTestMutation({testId: id}).unwrap();
        } else if (type === ResourceType.Transaction) {
          await deleteTransactionMutation({transactionId: id}).unwrap();
        }

        showNotification({
          type: 'success',
          title: `${capitalize(type)} deleted successfully`,
        });
        navigate('/');
      } catch (error) {
        showNotification({
          type: 'error',
          title: `Could not delete ${capitalize(type)}`,
          description: JSON.stringify(error),
        });
      }
    },
    [deleteTestMutation, deleteTransactionMutation, navigate, showNotification]
  );

  return useCallback(
    (id: string, name: string, type: ResourceType) => {
      onOpen({
        title: `Are you sure you want to delete “${name}”?`,
        async onConfirm() {
          await onConfirmDelete(id, type);
        },
      });
    },
    [onConfirmDelete, onOpen]
  );
};

export default useDeleteResource;
