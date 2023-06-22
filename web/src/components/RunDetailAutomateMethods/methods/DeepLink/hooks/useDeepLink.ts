import {useCallback, useState} from 'react';
import DeepLinkService, {TDeepLinkConfig} from 'services/DeepLink.service';

const useDeepLink = () => {
  const [deepLink, setDeepLink] = useState<string>('');

  const onGetDeepLink = useCallback((config: TDeepLinkConfig) => {
    const link = DeepLinkService.getLink(config);
    setDeepLink(link);
  }, []);

  return {deepLink, onGetDeepLink};
};

export default useDeepLink;
