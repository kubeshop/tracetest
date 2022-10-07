import {Typography} from 'antd';

import emptyIcon from 'assets/SpanAssertionsEmptyState.svg';
import * as S from './PaginatedList.styled';

const Empty = () => (
  <S.EmptyContainer>
    <img src={emptyIcon} />
    <Typography.Text disabled>No Data</Typography.Text>
  </S.EmptyContainer>
);

export default Empty;
