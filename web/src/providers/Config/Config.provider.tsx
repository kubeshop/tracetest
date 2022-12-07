import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo} from 'react';

import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {setUserPreference} from 'redux/slices/User.slice';
import {useGetConfigQuery} from 'redux/apis/TraceTest.api';
import {ConfigMode, TConfig} from 'types/Config.types';
import Config from 'models/Config.model';
import UserSelectors from '../../selectors/User.selectors';

interface IContext {
  config: TConfig;
  isLoading: boolean;
  isError: boolean;
  skipConfigSetup(): void;
  shouldDisplayConfigSetup: boolean;
}

const Context = createContext<IContext>({
  config: Config({}),
  skipConfigSetup: noop,
  isLoading: false,
  isError: false,
  shouldDisplayConfigSetup: false,
});

interface IProps {
  children: React.ReactNode;
}

export const useConfig = () => useContext(Context);

const ConfigProvider = ({children}: IProps) => {
  const dispatch = useAppDispatch();
  const {data: config = Config({}), isLoading, isError} = useGetConfigQuery({});
  const initConfigSetup = useAppSelector(state => UserSelectors.selectUserPreference(state, 'initConfigSetup'));

  const shouldDisplayConfigSetup = Boolean(initConfigSetup) && config.mode === ConfigMode.NO_TRACING_MODE;

  const skipConfigSetup = useCallback(() => {
    dispatch(
      setUserPreference({
        key: 'initConfigSetup',
        value: false,
      })
    );
  }, [dispatch]);

  const value = useMemo<IContext>(
    () => ({
      config,
      isLoading,
      isError,
      skipConfigSetup,
      shouldDisplayConfigSetup,
    }),
    [config, isError, isLoading, shouldDisplayConfigSetup, skipConfigSetup]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default ConfigProvider;
