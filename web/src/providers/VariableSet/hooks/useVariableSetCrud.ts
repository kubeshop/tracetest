import {useCallback} from 'react';
import TracetestAPI from 'redux/apis/Tracetest';
import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import {useNotification} from 'providers/Notification/Notification.provider';
import {TRawVariableSet} from 'models/VariableSet.model';

const {useCreateVariableSetMutation, useUpdateVariableSetMutation, useDeleteVariableSetMutation} =
  TracetestAPI.instance;

const useVariableSetCrud = () => {
  const [deleteVariableSet] = useDeleteVariableSetMutation();
  const [createVariableSet, {isLoading: isCreateLoading}] = useCreateVariableSetMutation();
  const [updateVariableSet, {isLoading: isUpdateLoading}] = useUpdateVariableSetMutation();
  const {onOpen} = useConfirmationModal();
  const {showNotification} = useNotification();

  const remove = useCallback(
    (id: string) => {
      onOpen({
        title: `Are you sure you want to delete the variable set?`,
        onConfirm: async () => {
          await deleteVariableSet({variableSetId: id});

          showNotification({
            type: 'success',
            title: 'Variable Set deleted successfully',
          });
        },
      });
    },
    [deleteVariableSet, onOpen, showNotification]
  );

  const edit = useCallback(
    async (variableSetId: string, variableSet: TRawVariableSet) => {
      await updateVariableSet({variableSetId, variableSet});
      showNotification({
        type: 'success',
        title: 'Variable Set updated successfully',
      });
    },
    [showNotification, updateVariableSet]
  );

  const create = useCallback(
    async (variableSet: TRawVariableSet) => {
      await createVariableSet(variableSet);

      showNotification({
        type: 'success',
        title: 'Variable Set created successfully',
      });
    },
    [createVariableSet, showNotification]
  );

  return {remove, isCreateLoading, isUpdateLoading, edit, create};
};

export default useVariableSetCrud;
