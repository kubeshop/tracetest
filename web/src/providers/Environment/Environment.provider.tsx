import {createContext, useContext, useMemo, useCallback} from 'react';
import {noop} from 'lodash';
import {TEnvironment} from 'types/Environment.types';
import {useGetEnvListQuery} from 'redux/apis/TraceTest.api';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {setUserPreference} from 'redux/slices/User.slice';
import EnvironmentSelectors from 'selectors/Environment.selectors';

interface IContext {
  environmentList: TEnvironment[];
  selectedEnvironment?: TEnvironment;
  setSelectedEnvironment(environment: TEnvironment): void;
}

export const Context = createContext<IContext>({
  environmentList: [],
  selectedEnvironment: undefined,
  setSelectedEnvironment: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useEnvironment = () => useContext(Context);

const EnvironmentProvider = ({children}: IProps) => {
  const {data: {items: environmentList = []} = {}} = useGetEnvListQuery({});
  const dispatch = useAppDispatch();
  const selectedEnvironment: TEnvironment | undefined = useAppSelector(EnvironmentSelectors.selectSelectedEnvironment);

  const setSelectedEnvironment = useCallback(
    (environment: TEnvironment) => {
      dispatch(
        setUserPreference({
          key: 'environmentId',
          value: environment.id,
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
    }),
    [environmentList, selectedEnvironment, setSelectedEnvironment]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default EnvironmentProvider;
