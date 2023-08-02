import {ReadOutlined} from '@ant-design/icons';
import {PropsWithChildren} from 'react';
import * as S from './DocsBanner.styled';

const DocsBanner: React.FC<PropsWithChildren<{}>> = ({children}) => {
  return (
    <S.DocsBannerContainer>
      <ReadOutlined />
      <S.Text>{children}</S.Text>
    </S.DocsBannerContainer>
  );
};

export default DocsBanner;
