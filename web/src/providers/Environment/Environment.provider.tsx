import {noop} from 'lodash';
import {createContext, useContext, useMemo, useCallback} from 'react';

import {useGetEnvironmentsQuery} from 'redux/apis/TraceTest.api';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {setUserPreference} from 'redux/slices/User.slice';
import EnvironmentSelectors from 'selectors/Environment.selectors';
import Environment from 'models/Environment.model';

interface IContext {
  environmentList: Environment[];
  selectedEnvironment?: Environment;
  setSelectedEnvironment(environment?: Environment): void;
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
  const dispatch = useAppDispatch();
  const {data: {items: environmentList = []} = {}, isLoading} = useGetEnvironmentsQuery({});
  const selectedEnvironment: Environment | undefined = useAppSelector(EnvironmentSelectors.selectSelectedEnvironment);

  const setSelectedEnvironment = useCallback(
    (environment?: Environment) => {
      dispatch(
        setUserPreference({
          key: 'environmentId',
          value: environment?.id || '',
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
    }),
    [environmentList, isLoading, selectedEnvironment, setSelectedEnvironment]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default EnvironmentProvider;
