import {Form} from 'antd';
import {useEffect, useMemo} from 'react';
import {useDataStore} from 'providers/DataStore/DataStore.provider';
import {TDraftDataStore} from 'types/DataStore.types';
import TracetestAPI from 'redux/apis/Tracetest/Tracetest.api';
import DataStoreService from 'services/DataStore.service';

export type TConnectionStatus = 'loading' | 'success' | 'error' | 'idle';

const useTestConnectionStatus = () => {
  // had to do this at this level because of this issue
  // https://github.com/reduxjs/redux-toolkit/issues/2055
  const {useResetTestOtlpConnectionMutation} = TracetestAPI.instance;
  const [resetOtlpCount] = useResetTestOtlpConnectionMutation();

  const {
    isTestConnectionLoading,
    isOtlpTestConnectionError,
    resetTestConnection,
    otlpTestConnectionResponse,
    onOtlpTestConnection,
  } = useDataStore();
  const form = Form.useFormInstance<TDraftDataStore>();

  const status = useMemo<TConnectionStatus>(() => {
    if (isTestConnectionLoading || !otlpTestConnectionResponse?.spanCount) return 'loading';
    if (!otlpTestConnectionResponse) return 'idle';
    if (isOtlpTestConnectionError) return 'error';

    return 'success';
  }, [isOtlpTestConnectionError, isTestConnectionLoading, otlpTestConnectionResponse]);

  // listens to all form changes
  const data = Form.useWatch([], form);
  const isOtlpBased = useMemo(() => DataStoreService.getIsOtlpBased(form.getFieldsValue()), [data]);

  useEffect(() => {
    if (isOtlpBased) onOtlpTestConnection();
  }, [isOtlpBased]);

  useEffect(() => {
    return () => {
      resetOtlpCount(undefined);
      resetTestConnection();
    };
  }, []);

  return {status, isOtlpBased, otlpResponse: otlpTestConnectionResponse, isLoading: isTestConnectionLoading};
};

export default useTestConnectionStatus;
