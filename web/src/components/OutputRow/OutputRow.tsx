import {Tag} from 'antd';

import {TTestOutput} from 'types/TestOutput.types';
import * as S from './OutputRow.styled';
import OutputRowActions from './OutputRowActions';
import OutputValue from './OutputValue';

interface IProps {
  index: number;
  output: TTestOutput;
  onDelete(index: number): void;
  onEdit(index: number): void;
}

const OutputRow = ({index, output, onEdit, onDelete}: IProps) => (
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
        <S.Entry>
          {!output.isDraft && Boolean(output.valueRun) && (
            <>
              <S.Key>Run value</S.Key>
              <OutputValue value={output.valueRun} />
            </>
          )}
        </S.Entry>
      </S.OutputDetails>
      <OutputRowActions name={output.name} onDelete={() => onDelete(index)} onEdit={() => onEdit(index)} />
    </S.Row>
  </S.Container>
);

export default OutputRow;
