import {capitalize} from 'lodash';
import * as S from './ImportSelector.styled';

interface IProps {
  onClick(): void;
  name: string;
  isSelected: boolean;
}

const ImportCard = ({name, onClick, isSelected}: IProps) => (
  <S.CardContainer data-cy={`${name.toLowerCase()}-plugin`} onClick={onClick}>
    <S.Circle>{isSelected && <S.Check />}</S.Circle>

    <S.CardContent>
      <S.CardTitle>Import via {capitalize(name)} </S.CardTitle>
    </S.CardContent>
  </S.CardContainer>
);

export default ImportCard;
