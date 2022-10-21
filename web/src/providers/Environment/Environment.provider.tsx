import {createContext, useContext, useMemo, useCallback, useEffect} from 'react';
import {noop} from 'lodash';
import {TEnvironment} from 'types/Environment.types';
import {useGetEnvListQuery, useLazyGetEnvironmentSecretListQuery} from 'redux/apis/TraceTest.api';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {setUserPreference} from 'redux/slices/User.slice';
import EnvironmentSelectors from 'selectors/Environment.selectors';

interface IContext {
  environmentList: TEnvironment[];
  selectedEnvironment?: TEnvironment;
  setSelectedEnvironment(environment?: TEnvironment): void;
  isLoading: boolean;
}

export const Context = createContext<IContext>({
  environmentList: [],
  selectedEnvironment: undefined,
  setSelectedEnvironment: noop,
  isLoading: true,
});

interface IProps {
  children: React.ReactNode;
}

export const useEnvironment = () => useContext(Context);

const EnvironmentProvider = ({children}: IProps) => {
  const {data: {items: environmentList = []} = {}, isLoading} = useGetEnvListQuery({});
  const [getEnvEntryList] = useLazyGetEnvironmentSecretListQuery({});
  const dispatch = useAppDispatch();
  const selectedEnvironment: TEnvironment | undefined = useAppSelector(EnvironmentSelectors.selectSelectedEnvironment);

  const getEnvironmentEntryList = useCallback(() => {
    if (selectedEnvironment) getEnvEntryList({environmentId: selectedEnvironment.id});
  }, [getEnvEntryList, selectedEnvironment]);

  const setSelectedEnvironment = useCallback(
    (environment?: TEnvironment) => {
      dispatch(
        setUserPreference({
          key: 'environmentId',
          value: environment?.id || '',
        })
      );
    },
    [dispatch]
  );

  useEffect(() => {
    getEnvironmentEntryList();
  }, [getEnvironmentEntryList]);

  const value = useMemo<IContext>(
    () => ({
      environmentList,
      selectedEnvironment,
      setSelectedEnvironment,
      isLoading,
    }),
    [environmentList, isLoading, selectedEnvironment, setSelectedEnvironment]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default EnvironmentProvider;
