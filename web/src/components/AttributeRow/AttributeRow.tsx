import {Popover} from 'antd';
import parse from 'html-react-parser';
import MarkdownIt from 'markdown-it';
import {useMemo} from 'react';

import AssertionResultChecks from 'components/AssertionResultChecks';
import AttributeActions from 'components/AttributeActions/AttributeActions';
import AttributeValue from 'components/AttributeValue';
import {OtelReference} from 'components/TestSpecForm/hooks/useGetOTELSemanticConventionAttributesInfo';
import TestOutput from 'models/TestOutput.model';
import TestRunOutput from 'models/TestRunOutput.model';
import SpanAttributeService from 'services/SpanAttribute.service';
import {TSpanFlatAttribute} from 'types/Span.types';
import {TTestSpecSummary} from 'types/TestRun.types';
import * as S from './AttributeRow.styled';

interface IProps {
  attribute: TSpanFlatAttribute;
  searchText?: string;
  semanticConventions: OtelReference;
  testSpecs?: TTestSpecSummary;
  testOutputs?: TestRunOutput[];
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
  onCreateOutput(attribute: TSpanFlatAttribute): void;
}

const AttributeRow = ({
  attribute: {key, value},
  attribute,
  searchText,
  semanticConventions,
  testSpecs,
  testOutputs,
  onCreateTestSpec,
  onCreateOutput,
}: IProps) => {
  const semanticConvention = SpanAttributeService.getReferencePicker(semanticConventions, key);
  const description = useMemo(() => parse(MarkdownIt().render(semanticConvention.description)), [semanticConvention]);
  const note = useMemo(() => parse(MarkdownIt().render(semanticConvention.note)), [semanticConvention]);
  const attributeTestSpecs = useMemo(
    () => SpanAttributeService.getAttributeTestSpecs(key, testSpecs),
    [key, testSpecs]
  );
  const attributeTestOutputs = useMemo(
    () => SpanAttributeService.getAttributeTestOutputs(key, testOutputs),
    [key, testOutputs]
  );

  const cypressKey = key.toLowerCase().replace('.', '-');

  const content = (
    <S.DetailContainer>
      {description}
      {note}
      <S.TagsContainer>
        {semanticConvention.tags.map(tag => (
          <S.Tag key={tag}>{tag}</S.Tag>
        ))}
      </S.TagsContainer>
    </S.DetailContainer>
  );

  return (
    <S.Container data-cy={`attribute-row-${cypressKey}`}>
      <S.Header>
        <S.SectionTitle>
          <S.AttributeTitle title={key} searchText={searchText} />

          {semanticConvention.description !== '' && (
            <Popover content={content} placement="right" title={<S.Title level={3}>{key}</S.Title>}>
              <S.InfoIcon />
            </Popover>
          )}

          {!!attributeTestOutputs.length && <S.OutputsMark outputs={attributeTestOutputs as TestOutput[]} />}
        </S.SectionTitle>

        <S.AttributeValueRow>
          <AttributeValue value={value} searchText={searchText} />
        </S.AttributeValueRow>
        <AssertionResultChecks failed={attributeTestSpecs.failed} passed={attributeTestSpecs.passed} />
      </S.Header>

      <AttributeActions attribute={attribute} onCreateTestOutput={onCreateOutput} onCreateTestSpec={onCreateTestSpec}>
        <a onClick={e => e.preventDefault()} style={{height: 'fit-content'}}>
          <S.MoreIcon />
        </a>
      </AttributeActions>
    </S.Container>
  );
};

export default AttributeRow;
