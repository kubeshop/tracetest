import {useMemo} from 'react';
import AssertionResultChecks from 'components/AssertionResultChecks';
import TestOutputMark from 'components/TestOutputMark';
import TestOutput from 'models/TestOutput.model';
import SpanAttributeService from 'services/SpanAttribute.service';
import * as S from './TestAttributeRow.styled';
import BaseAttributeRow from '../BaseAttributeRow/BaseAttributeRow';
import {IPropsAttributeRow} from '../SpanDetail';

const AttributeRow = ({
  attribute: {key},
  attribute,
  searchText,
  semanticConventions,
  testSpecs,
  testOutputs,
  onCreateTestSpec,
  onCreateOutput,
}: IPropsAttributeRow) => {
  const attributeTestSpecs = useMemo(
    () => SpanAttributeService.getAttributeTestSpecs(key, testSpecs),
    [key, testSpecs]
  );

  const attributeTestOutputs = useMemo(
    () => SpanAttributeService.getAttributeTestOutputs(key, testOutputs),
    [key, testOutputs]
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
        <AssertionResultChecks failed={attributeTestSpecs.failed} passed={attributeTestSpecs.passed} />
        {!!attributeTestOutputs.length && <TestOutputMark outputs={attributeTestOutputs as TestOutput[]} />}
      </S.Footer>
    </S.Container>
  );
};

export default AttributeRow;
