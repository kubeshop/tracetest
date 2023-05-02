import {noop} from 'lodash';
import {createContext, useContext, useMemo, useCallback, useState} from 'react';
import {useGetEnvironmentsQuery} from 'redux/apis/TraceTest.api';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {setUserPreference} from 'redux/slices/User.slice';
import EnvironmentSelectors from 'selectors/Environment.selectors';
import Environment from 'models/Environment.model';
import EnvironmentModal from 'components/EnvironmentModal';
import EnvironmentService from 'services/Environment.service';
import useEnvironmentCrud from './hooks/useEnvironmentCrud';

interface IContext {
  environmentList: Environment[];
  selectedEnvironment?: Environment;
  setSelectedEnvironment(environment?: Environment): void;
  isLoading: boolean;
  onOpenModal(draftEnvironment?: Environment): void;
  onDelete(id: string): void;
}

export const Context = createContext<IContext>({
  environmentList: [],
  selectedEnvironment: undefined,
  setSelectedEnvironment: noop,
  isLoading: true,
  onOpenModal: noop,
  onDelete: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useEnvironment = () => useContext(Context);

const EnvironmentProvider = ({children}: IProps) => {
  const dispatch = useAppDispatch();
  const {data: {items: environmentList = []} = {}, isLoading} = useGetEnvironmentsQuery({});
  const selectedEnvironment: Environment | undefined = useAppSelector(EnvironmentSelectors.selectSelectedEnvironment);
  const [environment, setEnvironment] = useState<Environment | undefined>(undefined);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const {remove, edit, create, isCreateLoading, isUpdateLoading} = useEnvironmentCrud();

  const onOpenModal = useCallback((draftEnvironment?: Environment) => {
    setIsModalOpen(true);
    setEnvironment(draftEnvironment);
  }, []);

  const onDelete = useCallback(
    (id: string) => {
      remove(id);
    },
    [remove]
  );

  const onSubmit = useCallback(
    (values: Environment) => {
      const request = EnvironmentService.getRequest(values);
      if (environment) {
        edit(environment.id, request);
      } else {
        create(request);
      }
      setIsModalOpen(false);
    },
    [create, edit, environment]
  );

  const setSelectedEnvironment = useCallback(
    (newEnvironment?: Environment) => {
      dispatch(
        setUserPreference({
          key: 'environmentId',
          value: newEnvironment?.id || '',
        })
      );
    },
    [dispatch]
  );

  const value = useMemo<IContext>(
    () => ({
      environmentList,
      selectedEnvironment,
      setSelectedEnvironment,
      isLoading,
      onOpenModal,
      onDelete,
    }),
    [environmentList, isLoading, onDelete, onOpenModal, selectedEnvironment, setSelectedEnvironment]
  );

  return (
    <>
      <Context.Provider value={value}>{children}</Context.Provider>
      <EnvironmentModal
        onSubmit={onSubmit}
        isLoading={isCreateLoading || isUpdateLoading}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        environment={environment}
      />
    </>
  );
};

export default EnvironmentProvider;
