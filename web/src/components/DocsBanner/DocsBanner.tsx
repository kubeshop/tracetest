import {ReadOutlined} from '@ant-design/icons';
import * as S from './DocsBanner.styled';

const DocsBanner: React.FC = ({children}) => {
  return (
    <S.DocsBannerContainer>
      <ReadOutlined />
      <S.Text>{children}</S.Text>
    </S.DocsBannerContainer>
  );
};

export default DocsBanner;
