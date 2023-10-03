import {useEffect} from 'react';
import {useParams} from 'react-router-dom';
import RouterMiddleware from '../redux/Router.middleware';

const useRouterSync = () => {
  const params = useParams();

  useEffect(() => {
    return RouterMiddleware.startListening({testId: params.testId, runId: params.runId});
  }, [params.runId, params.testId]);

  useEffect(() => {
    return RouterMiddleware.startListeningForLocationChange();
  }, []);
};

export default useRouterSync;
