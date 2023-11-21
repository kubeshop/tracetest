import {noop} from 'lodash';
import {createContext, useContext} from 'react';
import {NavigateFunction} from 'react-router-dom';

interface IContext {
  baseUrl: string;
  dashboardUrl: string;
  navigate: NavigateFunction;
}

export const Context = createContext<IContext>({
  baseUrl: '',
  dashboardUrl: '',
  navigate: noop,
});

export const useDashboard = () => useContext(Context);

interface IProps {
  children: React.ReactNode;
  value: IContext;
}

const DashboardProvider = ({children, value}: IProps) => {
  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default DashboardProvider;
