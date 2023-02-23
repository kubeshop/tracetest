import {useCallback} from 'react';
import {message} from 'antd';
import {
  useCreateEnvironmentMutation,
  useDeleteEnvironmentMutation,
  useUpdateEnvironmentMutation,
} from 'redux/apis/TraceTest.api';
import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import Environment from 'models/Environment.model';

const useEnvironmentCrud = () => {
  const [deleteEnvironment] = useDeleteEnvironmentMutation();
  const [createEnvironment, {isLoading: isCreateLoading}] = useCreateEnvironmentMutation();
  const [updateEnvironment, {isLoading: isUpdateLoading}] = useUpdateEnvironmentMutation();
  const {onOpen} = useConfirmationModal();

  const remove = useCallback(
    (id: string) => {
      onOpen({
        title: `Are you sure you want to delete the environment?`,
        onConfirm: () => deleteEnvironment({environmentId: id}),
      });
    },
    [deleteEnvironment, onOpen]
  );

  const edit = useCallback(
    async (environmentId: string, environment: Environment) => {
      await updateEnvironment({environmentId, environment});
      message.success('Environment updated successfully');
    },
    [updateEnvironment]
  );

  const create = useCallback(
    async (environment: Environment) => {
      await createEnvironment(environment);
      message.success('Environment created successfully');
    },
    [createEnvironment]
  );

  return {remove, isCreateLoading, isUpdateLoading, edit, create};
};

export default useEnvironmentCrud;
