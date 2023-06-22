import {useEffect} from 'react';
import AnalyticsService from 'services/Analytics/Analytics.service';

const withAnalytics = <P extends object>(Component: React.ComponentType<P>, name: string) => {
  const FunctionComponent = (props: P) => {
    useEffect(() => {
      AnalyticsService.page(name);
    }, [props]);

    return <Component {...props} />;
  };

  return FunctionComponent;
};

export default withAnalytics;
