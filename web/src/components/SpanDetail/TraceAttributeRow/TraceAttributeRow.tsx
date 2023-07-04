import {useMemo} from 'react';
import AnalyzerErrorsPopover from 'components/AnalyzerErrorsPopover';
import SpanAttributeService from 'services/SpanAttribute.service';
import * as S from './TraceAttributeRow.styled';
import BaseAttributeRow from '../BaseAttributeRow/BaseAttributeRow';
import {IPropsAttributeRow} from '../SpanDetail';

const AttributeRow = ({
  attribute: {key},
  attribute,
  searchText,
  semanticConventions,
  analyzerErrors,
  onCreateTestSpec,
  onCreateOutput,
}: IPropsAttributeRow) => {
  const attributeAnalyzerErrors = useMemo(
    () => SpanAttributeService.getAttributeAnalyzerErrors(key, analyzerErrors),
    [key, analyzerErrors]
  );

  return (
    <S.Container>
      <BaseAttributeRow
        attribute={attribute}
        searchText={searchText}
        semanticConventions={semanticConventions}
        onCreateTestSpec={onCreateTestSpec}
        onCreateOutput={onCreateOutput}
      />

      <S.Footer>
        {!!attributeAnalyzerErrors.length && <AnalyzerErrorsPopover errors={attributeAnalyzerErrors} />}
      </S.Footer>
    </S.Container>
  );
};

export default AttributeRow;
