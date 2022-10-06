import {Skeleton} from 'antd';

import * as S from './Home.styled';

const Loading = () => {
  return (
    <S.LoadingContainer direction="vertical">
      <Skeleton.Input active block size="small" />
      <Skeleton.Input active block size="small" />
      <Skeleton.Input active block size="small" />
    </S.LoadingContainer>
  );
};

export default Loading;
