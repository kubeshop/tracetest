import {PlusOutlined} from '@ant-design/icons';
import {Button} from 'antd';
import React, {useCallback, useEffect, useState} from 'react';
import {useNavigate} from 'react-router-dom';

import OutputModal from 'components/OutputModal/OutputModal';
import OutputRow from 'components/OutputRow';
import {toRawTestOutputs} from 'models/TestOutput.model';
import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {useSetTestOutputsMutation, useReRunMutation} from 'redux/apis/TraceTest.api';
import {selectTestOutputs, selectIsPending} from 'redux/testOutputs/selectors';
import {outputsInitiated, outputAdded, outputUpdated, outputDeleted, outputsReverted} from 'redux/testOutputs/slice';
import {useTest} from 'providers/Test/Test.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {TTestOutput} from 'types/TestOutput.types';
import * as S from './RunDetailTriggerResponse.styled';
import SkeletonResponse from './SkeletonResponse';

const ResponseOutputs = () => {
  const {
    run: {id: runId},
  } = useTestRun();
  const {
    test: {id: testId, outputs: testOutputs},
  } = useTest();

  const dispatch = useAppDispatch();
  const outputs = useAppSelector(state => selectTestOutputs(state, testId, runId));
  const isPending = useAppSelector(selectIsPending);
  const [setTestOutputs, {isLoading}] = useSetTestOutputsMutation();
  const [reRunTest, {isLoading: isLoadingReRub}] = useReRunMutation();
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const {onOpen} = useConfirmationModal();
  const navigate = useNavigate();

  useEffect(() => {
    dispatch(outputsInitiated(testOutputs || []));
  }, [dispatch, testOutputs]);

  const handleOnDelete = useCallback(
    (index: number) => {
      onOpen(`Are you sure you want to delete the output?`, () => {
        dispatch(outputDeleted(index));
      });
    },
    [dispatch, onOpen]
  );

  const handleOnEdit = useCallback((index: number) => {
    setSelectedIndex(index);
    setIsModalOpen(true);
  }, []);

  const handleOnSubmit = useCallback(
    (values: TTestOutput, isEditing: boolean) => {
      setIsModalOpen(false);
      if (isEditing) {
        dispatch(outputUpdated({index: selectedIndex, output: values}));
        return;
      }
      dispatch(outputAdded(values));
    },
    [dispatch, selectedIndex]
  );

  const handleOnCancel = useCallback(() => {
    dispatch(outputsReverted());
  }, [dispatch]);

  const handleOnSave = useCallback(async () => {
    await setTestOutputs({testId, testOutputs: toRawTestOutputs(outputs)}).unwrap();
    const run = await reRunTest({testId, runId}).unwrap();
    navigate(`/test/${testId}/run/${run.id}/trigger`);
  }, [navigate, outputs, reRunTest, runId, setTestOutputs, testId]);

  return !outputs ? (
    <SkeletonResponse />
  ) : (
    <>
      <S.HeadersList>
        {outputs.map((output, index) => (
          <OutputRow index={index} key={output.name} output={output} onDelete={handleOnDelete} onEdit={handleOnEdit} />
        ))}
      </S.HeadersList>

      <div>
        <Button
          icon={<PlusOutlined />}
          onClick={() => {
            setIsModalOpen(true);
            setSelectedIndex(-1);
          }}
          style={{fontWeight: 600, height: 'auto', padding: 0}}
          type="link"
        >
          Add Output
        </Button>
      </div>

      {isPending && (
        <S.Actions>
          <Button ghost loading={isLoading || isLoadingReRub} onClick={handleOnCancel} type="primary">
            Reset
          </Button>
          <Button loading={isLoading || isLoadingReRub} onClick={handleOnSave} type="primary">
            Save & Run
          </Button>
        </S.Actions>
      )}

      <OutputModal
        index={selectedIndex}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleOnSubmit}
        runId={runId}
        testId={testId}
      />
    </>
  );
};

export default ResponseOutputs;
