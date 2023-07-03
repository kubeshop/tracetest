import Highlighted from 'components/Highlighted';
import * as S from './BaseAttributeRow.styled';

interface IProps {
  searchText?: string;
  title: string;
}

const AttributeTitle = ({searchText = '', title}: IProps) => (
  <S.TextContainer>
    <S.Text type="secondary">
      <Highlighted text={title} highlight={searchText} />
    </S.Text>
  </S.TextContainer>
);

export default AttributeTitle;
