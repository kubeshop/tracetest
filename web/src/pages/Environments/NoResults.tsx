import {FC} from 'react';
import {Typography} from 'antd';
import * as S from './Envs.styled';

const NoResults: FC = () => {
  return (
    <S.NoResultsContainer>
      <S.NoResultsIcon />
      <S.NoResultsTitle>There are no test to show</S.NoResultsTitle>
      <Typography.Text>Start by creating a new test</Typography.Text>
    </S.NoResultsContainer>
  );
};

export default NoResults;
