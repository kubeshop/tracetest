import Highlighted from '../Highlighted';
import {Text, TextContainer} from './AttributeRow.styled';

interface IProps {
  searchText?: string;
  title: string;
}

const AttributeTitle = ({searchText = '', title}: IProps) => (
  <TextContainer>
    <Text type="secondary">
      <Highlighted text={title} highlight={searchText} />
    </Text>
  </TextContainer>
);

export default AttributeTitle;
