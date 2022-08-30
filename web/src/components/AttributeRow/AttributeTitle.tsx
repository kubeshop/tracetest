import {Tooltip} from 'antd';
import Highlighted from '../Highlighted';
import {Text, TextContainer} from './AttributeRow.styled';

interface IProps {
  title: string;
  searchText?: string;
}

export default ({searchText = '', title}: IProps): JSX.Element => {
  const textContainer = (
    <TextContainer>
      <Text type="secondary">
        <Highlighted text={title} highlight={searchText} />
      </Text>
    </TextContainer>
  );
  return title.length > 26 ? (
    <Tooltip title={title} arrowContent={null}>
      {textContainer}
    </Tooltip>
  ) : (
    textContainer
  );
};
