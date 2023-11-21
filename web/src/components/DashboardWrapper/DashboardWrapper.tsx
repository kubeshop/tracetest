import {useMemo} from 'react';
import {useNavigate} from 'react-router-dom';
import DashboardProvider from 'providers/Dashboard';
import {getServerBaseUrl} from 'utils/Common';

interface IProps {
  children: React.ReactNode;
}

const DashboardWrapper = ({children}: IProps) => {
  const navigate = useNavigate();
  const dashboardProviderValue = useMemo(() => ({baseUrl: '', dashboardUrl: getServerBaseUrl(), navigate}), [navigate]);

  return <DashboardProvider value={dashboardProviderValue}>{children}</DashboardProvider>;
};

export default DashboardWrapper;
