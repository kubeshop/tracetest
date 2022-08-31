import {Tooltip} from 'antd';
import Highlighted from '../Highlighted';
import {Text, TextContainer} from './AttributeRow.styled';

interface IProps {
  searchText?: string;
  title: string;
}

const AttributeTitle = ({searchText = '', title}: IProps) => {
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

export default AttributeTitle;
