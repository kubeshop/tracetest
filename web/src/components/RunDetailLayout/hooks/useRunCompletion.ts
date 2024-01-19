import {useCallback, useEffect, useState} from 'react';
import {isRunStateSucceeded} from 'models/TestRun.model';
import {useNotification} from 'providers/Notification/Notification.provider';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {ConfigMode} from 'types/DataStore.types';
import {useWizard} from 'providers/Wizard';

const useRunCompletion = () => {
  const {showNotification} = useNotification();
  const {run, isSkippedPolling} = useTestRun();
  const {dataStoreConfig} = useSettingsValues();
  const [prevState, setPrevState] = useState(run.state);
  const {onCompleteStep} = useWizard();

  const handleCompletion = useCallback(() => {
    const isNoTracingMode = dataStoreConfig.mode === ConfigMode.NO_TRACING_MODE;

    if (isRunStateSucceeded(run.state) && !isRunStateSucceeded(prevState)) {
      showNotification({
        type: 'success',
        title: isNoTracingMode
          ? 'Response received. Skipping looking for trace as you are in No-Tracing Mode'
          : 'Trace has been fetched successfully',
      });

      if (!isSkippedPolling) {
        onCompleteStep('create_test');
      }
    }

    setPrevState(run.state);
  }, [dataStoreConfig.mode, isSkippedPolling, onCompleteStep, prevState, run.state, showNotification]);

  useEffect(() => {
    handleCompletion();
  }, [dataStoreConfig.mode, handleCompletion, prevState, run.state, showNotification]);
};

export default useRunCompletion;
