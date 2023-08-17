import {noop} from 'lodash';
import {createContext, useContext, useMemo, useCallback, useState} from 'react';
import TracetestAPI from 'redux/apis/Tracetest';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {setUserPreference} from 'redux/slices/User.slice';
import VariableSetSelectors from 'selectors/VariableSet.selectors';
import VariableSet from 'models/VariableSet.model';
import VariableSetModal from 'components/VariableSetModal';
import VariableSetService from 'services/VariableSet.service';
import useVariableSetCrud from './hooks/useVariableSetCrud';

const {useGetVariableSetsQuery} = TracetestAPI.instance;

interface IContext {
  variableSetList: VariableSet[];
  selectedVariableSet?: VariableSet;
  setSelectedVariableSet(variableSet?: VariableSet): void;
  isLoading: boolean;
  onOpenModal(draftVariableSet?: VariableSet): void;
  onDelete(id: string): void;
}

export const Context = createContext<IContext>({
  variableSetList: [],
  setSelectedVariableSet: noop,
  isLoading: true,
  onOpenModal: noop,
  onDelete: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useVariableSet = () => useContext(Context);

const VariableSetProvider = ({children}: IProps) => {
  const dispatch = useAppDispatch();
  const {data: {items: variableSetList = []} = {}, isLoading} = useGetVariableSetsQuery({});
  const selectedVariableSet: VariableSet | undefined = useAppSelector(VariableSetSelectors.selectSelectedVariableSet);
  const [variableSet, setVariableSet] = useState<VariableSet | undefined>(undefined);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const {remove, edit, create, isCreateLoading, isUpdateLoading} = useVariableSetCrud();

  const onOpenModal = useCallback((draftVariableSet?: VariableSet) => {
    setIsModalOpen(true);
    setVariableSet(draftVariableSet);
  }, []);

  const onDelete = useCallback(
    (id: string) => {
      remove(id);
    },
    [remove]
  );

  const onSubmit = useCallback(
    (values: VariableSet) => {
      const request = VariableSetService.getRequest(values);
      if (variableSet) {
        edit(variableSet.id, request);
      } else {
        create(request);
      }
      setIsModalOpen(false);
    },
    [create, edit, variableSet]
  );

  const setSelectedVariableSet = useCallback(
    (newVariableSet?: VariableSet) => {
      dispatch(
        setUserPreference({
          key: 'variableSetId',
          value: newVariableSet?.id || '',
        })
      );
    },
    [dispatch]
  );

  const value = useMemo<IContext>(
    () => ({
      variableSetList,
      selectedVariableSet,
      setSelectedVariableSet,
      isLoading,
      onOpenModal,
      onDelete,
    }),
    [variableSetList, selectedVariableSet, setSelectedVariableSet, isLoading, onOpenModal, onDelete]
  );

  return (
    <>
      <Context.Provider value={value}>{children}</Context.Provider>
      <VariableSetModal
        onSubmit={onSubmit}
        isLoading={isCreateLoading || isUpdateLoading}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        variableSet={variableSet}
      />
    </>
  );
};

export default VariableSetProvider;
