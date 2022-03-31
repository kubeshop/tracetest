import {FC} from 'react';
import {Typography} from 'antd';
import * as S from './Home.styled';

const NoResults: FC = () => {
  return (
    <S.NoResultsContainer>
      <S.NoResultsIcon />
      <S.NoResultsTitle>The are no test to show</S.NoResultsTitle>
      <Typography.Text>Please create the new test</Typography.Text>
    </S.NoResultsContainer>
  );
};

export default NoResults;
