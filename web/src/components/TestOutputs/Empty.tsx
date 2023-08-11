import {ADD_TEST_OUTPUTS_DOCUMENTATION_URL} from 'constants/Common.constants';
import * as S from './TestOutputs.styled';

const Empty = () => (
  <S.EmptyContainer data-cy="empty-test-outputs">
    <S.EmptyIcon />
    <S.EmptyTitle>There are no outputs for this test</S.EmptyTitle>
    <S.EmptyText>Outputs create variables that are used in tests within a test suite.</S.EmptyText>
    <S.EmptyText>
      Learn more about outputs{' '}
      <a href={ADD_TEST_OUTPUTS_DOCUMENTATION_URL} target="_blank">
        here
      </a>
    </S.EmptyText>
  </S.EmptyContainer>
);

export default Empty;
