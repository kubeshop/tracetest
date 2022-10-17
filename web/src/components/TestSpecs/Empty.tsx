import {ADD_TEST_SPECS_DOCUMENTATION_URL} from 'constants/Common.constants';
import * as S from './TestSpecs.styled';

const Empty = () => (
  <S.EmptyContainer data-cy="empty-test-specs">
    <S.EmptyIcon />
    <S.EmptyTitle>There are no specs for this test</S.EmptyTitle>
    <S.EmptyText>Add a Test Spec to validate your trace.</S.EmptyText>
    <S.EmptyText>
      Learn more about writing specs{' '}
      <a href={ADD_TEST_SPECS_DOCUMENTATION_URL} target="_blank">
        here
      </a>
    </S.EmptyText>
  </S.EmptyContainer>
);

export default Empty;
