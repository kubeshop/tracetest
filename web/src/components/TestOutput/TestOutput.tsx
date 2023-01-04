import {Tag} from 'antd';

import AttributeValue from 'components/AttributeValue';
import {TTestOutput} from 'types/TestOutput.types';
import * as S from './TestOutput.styled';
import Actions from './Actions';

interface IProps {
  index: number;
  output: TTestOutput;
  onDelete(index: number): void;
  onEdit(values: TTestOutput): void;
}

const TestOutput = ({index, output, onEdit, onDelete}: IProps) => (
  <S.Container $isDeleted={output.isDeleted} data-cy="output-item-container">
    {output.isDraft && (
      <S.Row $justifyContent="flex-end">
        <Tag data-cy="output-pending-tag">pending {output.isDeleted && '/ deleted'}</Tag>
      </S.Row>
    )}
    <S.Row>
      <S.OutputDetails>
        <S.Entry>
          <S.Key>Name</S.Key>
          <S.Value>{output.name}</S.Value>
        </S.Entry>
        <S.Entry>
          <S.Key>Selector</S.Key>
          <S.Value>{output.selector}</S.Value>
        </S.Entry>
        <S.Entry>
          <S.Key>Value</S.Key>
          <S.Value>{output.value}</S.Value>
        </S.Entry>
      </S.OutputDetails>
      <Actions isDeleted={output.isDeleted} onDelete={() => onDelete(index)} onEdit={() => onEdit(output)} />
    </S.Row>
    <S.Row>
      <S.Entry>
        {!output.isDraft && Boolean(output.valueRun) && (
          <>
            <S.Key>Run value</S.Key>
            <AttributeValue value={output.valueRun} />
          </>
        )}
      </S.Entry>
    </S.Row>
  </S.Container>
);

export default TestOutput;
