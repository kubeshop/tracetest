import {MoreOutlined} from '@ant-design/icons';
import {Dropdown, Menu, message, Popover} from 'antd';
import parse from 'html-react-parser';
import MarkdownIt from 'markdown-it';
import {useMemo} from 'react';

import AttributeValue from 'components/AttributeValue';
import {OtelReference} from 'components/TestSpecForm/hooks/useGetOTELSemanticConventionAttributesInfo';
import SpanAttributeService from 'services/SpanAttribute.service';
import {TResultAssertions} from 'types/Assertion.types';
import {TSpanFlatAttribute} from 'types/Span.types';
import {TTestOutput} from 'types/TestOutput.types';
import * as S from './AttributeRow.styled';
import AssertionResultChecks from '../AssertionResultChecks/AssertionResultChecks';

interface IProps {
  assertions?: TResultAssertions;
  attribute: TSpanFlatAttribute;
  searchText?: string;
  onCopy(value: string): void;
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
  onCreateOutput(attribute: TSpanFlatAttribute): void;
  semanticConventions: OtelReference;
  outputs: TTestOutput[];
}

enum Action {
  Copy = '0',
  Create_Spec = '1',
  Create_Output = '2',
}

const AttributeRow = ({
  assertions = {},
  attribute: {key, value},
  attribute,
  onCopy,
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

  const handleOnClick = ({key: option}: {key: string}) => {
    if (option === Action.Copy) {
      message.success('Value copied to the clipboard');
      return onCopy(value);
    }

    if (option === Action.Create_Spec) {
      return onCreateTestSpec(attribute);
    }

    if (option === Action.Create_Output) {
      return onCreateOutput(attribute);
    }
  };

  const cypressKey = key.toLowerCase().replace('.', '-');

  const menu = (
    <Menu
      items={[
        {
          label: 'Copy value',
          key: Action.Copy,
        },
        {
          label: 'Create output',
          key: Action.Create_Output,
        },
        {
          label: 'Create test spec',
          key: Action.Create_Spec,
        },
      ]}
      onClick={handleOnClick}
    />
  );

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

      <Dropdown overlay={menu}>
        <a onClick={e => e.preventDefault()}>
          <MoreOutlined />
        </a>
      </Dropdown>
    </S.Container>
  );
};

export default AttributeRow;
