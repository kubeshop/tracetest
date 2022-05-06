import {CloseCircleFilled} from '@ant-design/icons';
import {Button, Typography} from 'antd';
import React from 'react';
import {ITestRunResult} from 'types/TestRunResult.types';
import {useRunTestMutation} from 'redux/apis/Test.api';
import * as S from '../Trace.styled';

interface IProps {
  displayError: boolean;
  onDismissTrace: () => void;
  testId: string;

  onRunTest(result: ITestRunResult): void;
}

export const TraceError = ({displayError, onDismissTrace, onRunTest, testId}: IProps): JSX.Element | null => {
  const [runNewTest] = useRunTestMutation();
  return displayError ? (
    <S.FailedTrace>
      <CloseCircleFilled style={{color: 'red', fontSize: 32}} />
      <Typography.Title level={2}>Test Run Failed</Typography.Title>
      <div style={{display: 'grid', gap: 8, gridTemplateColumns: '1fr 1fr'}}>
        <Button
          onClick={async () => {
            const result = await runNewTest(testId).unwrap();
            onRunTest(result);
          }}
        >
          Rerun Test
        </Button>
        <Button onClick={onDismissTrace}>Cancel</Button>
      </div>
    </S.FailedTrace>
  ) : null;
};
