import {RESOURCE_SEMANTIC_CONVENTIONS_URL, TRACE_SEMANTIC_CONVENTIONS_URL} from 'constants/Common.constants';
import * as S from './AttributeList.styled';

const EmptyAttributeList: React.FC = () => {
  return (
    <S.EmptyAttributeList data-cy="empty-attribute-list">
      <S.EmptyIcon />
      <S.EmptyTitle>Looking for attributes?</S.EmptyTitle>
      <S.EmptyText>
        Take a look a the open telemetry specification for{' '}
        <a href={TRACE_SEMANTIC_CONVENTIONS_URL} target="_blank">
          trace semantic conventions
        </a>{' '}
        and here for{' '}
        <a href={RESOURCE_SEMANTIC_CONVENTIONS_URL} target="_blank">
          resource semantic conventions
        </a>
        .
      </S.EmptyText>
    </S.EmptyAttributeList>
  );
};
export default EmptyAttributeList;
