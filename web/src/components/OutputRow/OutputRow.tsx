import {capitalize} from 'lodash';
import {TOutput} from 'types/Output.types';
import * as S from './OutputRow.styled';
import OutputRowActions from './OutputRowActions';

interface IProps {
  output: TOutput;
  onDelete(id: string): void;
  onEdit(id: TOutput): void;
}

const OutputRow = ({output: {id, source, attribute, regex, selector}, output, onEdit, onDelete}: IProps) => {
  return (
    <S.Container>
      <S.OutputDetails>
        <S.Entry>
          <S.Key>Source</S.Key>
          <S.Value>{capitalize(source)}</S.Value>
        </S.Entry>
        <S.Entry>
          <S.Key>Attribute</S.Key>
          <S.Value>{attribute}</S.Value>
        </S.Entry>
        {selector && (
          <S.Entry>
            <S.Key>Selector</S.Key>
            <S.Value>{selector}</S.Value>
          </S.Entry>
        )}
        {regex && (
          <S.Entry>
            <S.Key>Regex</S.Key>
            <S.Value>{regex}</S.Value>
          </S.Entry>
        )}
      </S.OutputDetails>
      <OutputRowActions outputId={id} onDelete={() => onDelete(id)} onEdit={() => onEdit(output)} />
    </S.Container>
  );
};

export default OutputRow;
