import {PlusOutlined} from '@ant-design/icons';
import {Button} from 'antd';

import OutputRow from 'components/OutputRow';
import {selectTestOutputs, selectIsPending} from 'redux/testOutputs/selectors';
import {useTest} from 'providers/Test/Test.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useAppSelector} from 'redux/hooks';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import * as S from './RunDetailTriggerResponse.styled';
import SkeletonResponse from './SkeletonResponse';

const ResponseOutputs = () => {
  const {
    run: {id: runId},
  } = useTestRun();
  const {
    test: {id: testId},
  } = useTest();
  const outputs = useAppSelector(state => selectTestOutputs(state, testId, runId));
  const isPending = useAppSelector(selectIsPending);
  const {onModalOpenFromIndex, onDelete, onCancel, onSave, isLoading, onModalOpen} = useTestOutput();

  return !outputs ? (
    <SkeletonResponse />
  ) : (
    <>
      <S.HeadersList>
        {outputs.map((output, index) => (
          <OutputRow
            index={index}
            key={output.name}
            output={output}
            onDelete={onDelete}
            onEdit={onModalOpenFromIndex}
          />
        ))}
      </S.HeadersList>

      <div>
        <Button
          data-cy="output-add-button"
          icon={<PlusOutlined />}
          onClick={() => onModalOpen()}
          style={{fontWeight: 600, height: 'auto', padding: 0}}
          type="link"
        >
          Add Output
        </Button>
      </div>

      {isPending && (
        <S.Actions>
          <Button data-cy="output-reset-button" ghost loading={isLoading} onClick={onCancel} type="primary">
            Reset
          </Button>
          <Button data-cy="output-publish-button" loading={isLoading} onClick={() => onSave(outputs)} type="primary">
            Save & Run
          </Button>
        </S.Actions>
      )}
    </>
  );
};

export default ResponseOutputs;
