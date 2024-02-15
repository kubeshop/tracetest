import {useCallback} from 'react';
import {capitalize} from 'lodash';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import TracetestAPI from 'redux/apis/Tracetest';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {ResourceType} from 'types/Resource.type';
import {useNotification} from 'providers/Notification/Notification.provider';

const {useDeleteTestByIdMutation, useDeleteTestSuiteByIdMutation} = TracetestAPI.instance;

const useDeleteResource = () => {
  const [deleteTestMutation] = useDeleteTestByIdMutation();
  const [deleteTestSuiteMutation] = useDeleteTestSuiteByIdMutation();
  const {navigate} = useDashboard();
  const {showNotification} = useNotification();

  const {onOpen} = useConfirmationModal();

  const onConfirmDelete = useCallback(
    async (id: string, type: ResourceType) => {
      try {
        let path = '/tests';
        if (type === ResourceType.Test) {
          TestAnalyticsService.onDeleteTest();
          await deleteTestMutation({testId: id}).unwrap();
        } else if (type === ResourceType.TestSuite) {
          path = '/testsuites';
          await deleteTestSuiteMutation({testSuiteId: id}).unwrap();
        }

        showNotification({
          type: 'success',
          title: `${capitalize(type)} deleted successfully`,
        });

        navigate(path);
      } catch (error) {
        showNotification({
          type: 'error',
          title: `Could not delete ${capitalize(type)}`,
          description: JSON.stringify(error),
        });
      }
    },
    [deleteTestMutation, deleteTestSuiteMutation, navigate, showNotification]
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
