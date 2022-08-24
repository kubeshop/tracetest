import {useEffect} from 'react';
import {useParams} from 'react-router-dom';
import RouterMiddleware from '../redux/Router.middleware';

const useRouterSync = () => {
  const params = useParams();

  useEffect(() => {
    return RouterMiddleware.startListening(params);
  }, [params]);
};

export default useRouterSync;
