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
  const {useLazyTestOtlpConnectionQuery, useResetTestOtlpConnectionMutation} = TracetestAPI.instance;
  const [resetOtlpCount] = useResetTestOtlpConnectionMutation();

  const {isTestConnectionLoading, resetTestConnection} = useDataStore();
  const form = Form.useFormInstance<TDraftDataStore>();
  const [testOtlpConnection, {isLoading, data: otlpResponse, isError}] = useLazyTestOtlpConnectionQuery({
    pollingInterval: 1000,
  });

  const status = useMemo<TConnectionStatus>(() => {
    if (isTestConnectionLoading || !otlpResponse?.spanCount) return 'loading';
    if (!otlpResponse) return 'idle';
    if (isError) return 'error';

    return 'success';
  }, [isError, isTestConnectionLoading, otlpResponse]);

  // listens to all form changes
  const data = Form.useWatch([], form);
  const isOtlpBased = useMemo(() => DataStoreService.getIsOtlpBased(form.getFieldsValue()), [data]);

  useEffect(() => {
    if (isOtlpBased) testOtlpConnection(undefined);
  }, [isOtlpBased]);

  useEffect(() => {
    return () => {
      resetOtlpCount(undefined);
      resetTestConnection();
    };
  }, []);

  return {status, isOtlpBased, otlpResponse, isLoading: isLoading || isTestConnectionLoading};
};

export default useTestConnectionStatus;
