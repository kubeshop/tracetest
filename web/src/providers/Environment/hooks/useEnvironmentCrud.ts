import {useCallback} from 'react';
import {
  useCreateEnvironmentMutation,
  useDeleteEnvironmentMutation,
  useUpdateEnvironmentMutation,
} from 'redux/apis/Tracetest';
import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import {useNotification} from 'providers/Notification/Notification.provider';
import {TRawEnvironment} from 'models/Environment.model';

const useEnvironmentCrud = () => {
  const [deleteEnvironment] = useDeleteEnvironmentMutation();
  const [createEnvironment, {isLoading: isCreateLoading}] = useCreateEnvironmentMutation();
  const [updateEnvironment, {isLoading: isUpdateLoading}] = useUpdateEnvironmentMutation();
  const {onOpen} = useConfirmationModal();
  const {showNotification} = useNotification();

  const remove = useCallback(
    (id: string) => {
      onOpen({
        title: `Are you sure you want to delete the environment?`,
        onConfirm: async () => {
          await deleteEnvironment({environmentId: id});

          showNotification({
            type: 'success',
            title: 'Environment deleted successfully',
          });
        },
      });
    },
    [deleteEnvironment, onOpen, showNotification]
  );

  const edit = useCallback(
    async (environmentId: string, environment: TRawEnvironment) => {
      await updateEnvironment({environmentId, environment});
      showNotification({
        type: 'success',
        title: 'Environment updated successfully',
      });
    },
    [showNotification, updateEnvironment]
  );

  const create = useCallback(
    async (environment: TRawEnvironment) => {
      await createEnvironment(environment);

      showNotification({
        type: 'success',
        title: 'Environment created successfully',
      });
    },
    [createEnvironment, showNotification]
  );

  return {remove, isCreateLoading, isUpdateLoading, edit, create};
};

export default useEnvironmentCrud;
