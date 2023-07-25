import {useMemo} from 'react';
import {useNavigate} from 'react-router-dom';
import DashboardProvider from 'providers/Dashboard';

interface IProps {
  children: React.ReactNode;
}

const DashboardWrapper = ({children}: IProps) => {
  const navigate = useNavigate();
  const dashboardProviderValue = useMemo(() => ({baseUrl: '', navigate}), [navigate]);

  return <DashboardProvider value={dashboardProviderValue}>{children}</DashboardProvider>;
};

export default DashboardWrapper;
