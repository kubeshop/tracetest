import * as S from './TraceAttributeRow.styled';
import BaseAttributeRow from '../BaseAttributeRow/BaseAttributeRow';
import {IPropsAttributeRow} from '../SpanDetail';

const AttributeRow = ({
  attribute,
  searchText,
  semanticConventions,
  onCreateTestSpec,
  onCreateOutput,
}: IPropsAttributeRow) => (
  <S.Container>
    <BaseAttributeRow
      attribute={attribute}
      searchText={searchText}
      semanticConventions={semanticConventions}
      onCreateTestSpec={onCreateTestSpec}
      onCreateOutput={onCreateOutput}
    />
  </S.Container>
);

export default AttributeRow;
