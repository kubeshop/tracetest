import {useCallback, useState} from 'react';
import DeepLinkService, {TDeepLinkConfig} from 'services/DeepLink.service';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';

const useDeepLink = () => {
  const [deepLink, setDeepLink] = useState<string>('');
  const {dashboardUrl} = useDashboard();

  const onGetDeepLink = useCallback(
    (config: TDeepLinkConfig) => {
      const link = DeepLinkService.getLink({...config, baseUrl: dashboardUrl});
      setDeepLink(link);
    },
    [dashboardUrl]
  );

  return {deepLink, onGetDeepLink};
};

export default useDeepLink;
