import {Tag, Tooltip} from 'antd';
import {useCallback} from 'react';
import AttributeValue from 'components/AttributeValue';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import {useSpan} from 'providers/Span/Span.provider';
import TestOutputModel from 'models/TestOutput.model';
import {selectIsSelectedOutput} from 'redux/testOutputs/selectors';
import {useAppSelector} from 'redux/hooks';
import * as S from './TestOutput.styled';
import Actions from './Actions';

interface IProps {
  index: number;
  output: TestOutputModel;
  onDelete(index: number): void;
  onEdit(values: TestOutputModel): void;
}

const TestOutput = ({
  index,
  output: {id, name, isDeleted, isDraft, spanId, selector, value, valueRun, valueRunDraft, error},
  output,
  onEdit,
  onDelete,
}: IProps) => {
  const {onSelectedOutputs} = useTestOutput();
  const {onSelectSpan} = useSpan();
  const isSelected = useAppSelector(state => selectIsSelectedOutput(state, id));

  const handleOutputClick = useCallback(() => {
    onSelectedOutputs([output]);
    onSelectSpan(spanId);
  }, [onSelectSpan, onSelectedOutputs, output, spanId]);

  return (
    <S.Container
      $isDeleted={isDeleted}
      data-cy="output-item-container"
      $isSelected={isSelected}
      onClick={handleOutputClick}
    >
      {isDraft && (
        <S.Row $justifyContent="flex-end">
          <Tag data-cy="output-pending-tag">pending {isDeleted && '/ deleted'}</Tag>
        </S.Row>
      )}
      <S.Row>
        <S.OutputDetails>
          <S.Entry>
            <S.Key>Name</S.Key>
            <S.Value>{name}</S.Value>
          </S.Entry>
          <S.Entry>
            <S.Key>Selector</S.Key>
            <S.Value>{selector}</S.Value>
          </S.Entry>
          <S.Entry>
            <S.Key>Value</S.Key>
            <S.Value>{value}</S.Value>
          </S.Entry>
        </S.OutputDetails>
        <Actions isDeleted={isDeleted} onDelete={() => onDelete(index)} onEdit={() => onEdit(output)} />
      </S.Row>
      <S.Row>
        <S.Entry>
          {!isDraft && Boolean(valueRun) && (
            <>
              <S.Key>Run value</S.Key>
              <AttributeValue value={valueRun} />
            </>
          )}
          {isDraft && Boolean(valueRunDraft) && (
            <>
              <S.Key>Run value</S.Key>
              <AttributeValue value={valueRunDraft} />
            </>
          )}
          {error && (
            <Tooltip title={error}>
              <S.IconWarning />
            </Tooltip>
          )}
        </S.Entry>
      </S.Row>
    </S.Container>
  );
};

export default TestOutput;
