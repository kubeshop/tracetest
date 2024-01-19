import {Form} from 'antd';
import {useCallback, useEffect, useMemo, useState} from 'react';
import {useDataStore} from 'providers/DataStore/DataStore.provider';
import {TDraftDataStore} from 'types/DataStore.types';
import TracetestAPI from 'redux/apis/Tracetest/Tracetest.api';
import DataStoreService from 'services/DataStore.service';

export type TConnectionStatus = 'loading' | 'success' | 'error' | 'idle';

const POLL_INTERVAL = 1000;

const useTestConnectionStatus = () => {
  // had to do this at this level because of this issue
  // https://github.com/reduxjs/redux-toolkit/issues/2055
  const {useResetTestOtlpConnectionMutation, useLazyTestOtlpConnectionQuery} = TracetestAPI.instance;
  const [resetOtlpCount] = useResetTestOtlpConnectionMutation();
  const [pollingInterval, setPollInterval] = useState<number | undefined>(undefined);
  const [
    testOtlpConnection,
    {isLoading: isOtlpTestConnectionLoading, data: otlpTestConnectionResponse, isError: isOtlpTestConnectionError},
  ] = useLazyTestOtlpConnectionQuery({
    pollingInterval,
  });

  const {isTestConnectionLoading, resetTestConnection, onSetOtlpTestConnectionResponse} = useDataStore();
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

  const onStartPolling = useCallback(async () => {
    setPollInterval(POLL_INTERVAL);
    await testOtlpConnection(undefined).unwrap();
  }, [testOtlpConnection]);

  const onReset = useCallback(() => {
    // stops the polling
    setPollInterval(undefined);

    // resets backend otlp count to 0
    resetOtlpCount(undefined);

    // resets the test connection result to undefined
    onSetOtlpTestConnectionResponse(undefined);

    // resets the direct test connection
    resetTestConnection();
  }, [onSetOtlpTestConnectionResponse, resetOtlpCount, resetTestConnection]);

  useEffect(() => {
    // resets the test connection results and data
    onReset();

    // if its otlp, starts the polling mechanism
    if (isOtlpBased) {
      onStartPolling();
    }
  }, [isOtlpBased]);

  useEffect(() => {
    /// if we are polling, refresh the provider result with the response information
    if (pollingInterval) onSetOtlpTestConnectionResponse(otlpTestConnectionResponse);
  }, [onSetOtlpTestConnectionResponse, otlpTestConnectionResponse, pollingInterval]);

  return {
    status,
    isOtlpBased,
    otlpResponse: otlpTestConnectionResponse,
    isLoading: isTestConnectionLoading || isOtlpTestConnectionLoading,
  };
};

export default useTestConnectionStatus;
