import {Popover} from 'antd';
import parse from 'html-react-parser';
import MarkdownIt from 'markdown-it';
import {useMemo} from 'react';

import AttributeActions from 'components/AttributeActions/AttributeActions';
import AttributeValue from 'components/AttributeValue';
import {OtelReference} from 'components/TestSpecForm/hooks/useGetOTELSemanticConventionAttributesInfo';
import SpanAttributeService from 'services/SpanAttribute.service';
import {TResultAssertions} from 'types/Assertion.types';
import {TSpanFlatAttribute} from 'types/Span.types';
import TestOutput from 'models/TestOutput.model';
import * as S from './AttributeRow.styled';
import AssertionResultChecks from '../AssertionResultChecks/AssertionResultChecks';

interface IProps {
  assertions?: TResultAssertions;
  attribute: TSpanFlatAttribute;
  searchText?: string;
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
  onCreateOutput(attribute: TSpanFlatAttribute): void;
  semanticConventions: OtelReference;
  outputs: TestOutput[];
}

const AttributeRow = ({
  assertions = {},
  attribute: {key, value},
  attribute,
  onCreateTestSpec,
  searchText,
  semanticConventions,
  onCreateOutput,
  outputs,
}: IProps) => {
  const semanticConvention = SpanAttributeService.getReferencePicker(semanticConventions, key);
  const description = useMemo(() => parse(MarkdownIt().render(semanticConvention.description)), [semanticConvention]);
  const note = useMemo(() => parse(MarkdownIt().render(semanticConvention.note)), [semanticConvention]);
  const {failed, passed} = useMemo(
    () => SpanAttributeService.getAttributeAssertionResults(key, assertions),
    [assertions, key]
  );
  const attributeOutputs = useMemo(
    () => SpanAttributeService.getOutputsFromAttributeName(key, outputs),
    [key, outputs]
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

          {!!attributeOutputs.length && <S.OutputsMark outputs={attributeOutputs} />}
        </S.SectionTitle>

        <S.AttributeValueRow>
          <AttributeValue value={value} searchText={searchText} />
        </S.AttributeValueRow>
        <AssertionResultChecks failed={failed} passed={passed} />
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
