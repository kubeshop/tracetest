import {MoreOutlined} from '@ant-design/icons';
import {Dropdown, Menu, message, Popover} from 'antd';
import parse from 'html-react-parser';
import MarkdownIt from 'markdown-it';
import React, {useMemo} from 'react';

import AttributeValue from 'components/AttributeValue';
import {OtelReference} from 'components/TestSpecForm/hooks/useGetOTELSemanticConventionAttributesInfo';
import SpanAttributeService from 'services/SpanAttribute.service';
import {IResult} from 'types/Assertion.types';
import {TSpanFlatAttribute} from 'types/Span.types';
import AttributeCheck from './AttributeCheck';
import * as S from './AttributeRow.styled';

interface IProps {
  assertionsFailed?: IResult[];
  assertionsPassed?: IResult[];
  attribute: TSpanFlatAttribute;
  searchText?: string;
  onCopy(value: string): void;
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
  semanticConventions: OtelReference;
}

enum Action {
  Copy = '0',
  Create_Spec = '1',
}

const AttributeRow = ({
  assertionsFailed,
  assertionsPassed,
  attribute: {key, value},
  attribute,
  onCopy,
  onCreateTestSpec,
  searchText,
  semanticConventions,
}: IProps) => {
  const passedCount = assertionsPassed?.length ?? 0;
  const failedCount = assertionsFailed?.length ?? 0;
  const semanticConvention = SpanAttributeService.getReferencePicker(semanticConventions, key);
  const description = useMemo(() => parse(MarkdownIt().render(semanticConvention.description)), [semanticConvention]);

  const handleOnClick = ({key: option}: {key: string}) => {
    if (option === Action.Copy) {
      message.success('Value copied to the clipboard');
      return onCopy(value);
    }

    if (option === Action.Create_Spec) {
      return onCreateTestSpec(attribute);
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
          label: 'Create test spec',
          key: Action.Create_Spec,
        },
      ]}
      onClick={handleOnClick}
    />
  );

  const content = (
    <S.DetailContainer>
      {semanticConvention.description !== '' ? description : 'We have not found a description for this attribute'}
      <S.TagsContainer>
        {semanticConvention.tags.map(tag => (
          <S.Tag key={tag}>{tag}</S.Tag>
        ))}
      </S.TagsContainer>
    </S.DetailContainer>
  );

  return (
    <S.Container data-cy={`attribute-row-${cypressKey}`}>
      <Popover content={content} placement="right" title={<S.Title level={3}>{key}</S.Title>} trigger="click">
        <S.Header>
          <S.AttributeTitle title={key} searchText={searchText} />

          <S.AttributeValueRow>
            <AttributeValue value={value} searchText={searchText} />
          </S.AttributeValueRow>
          {passedCount > 0 && <AttributeCheck items={assertionsPassed!} type="success" />}
          {failedCount > 0 && <AttributeCheck items={assertionsFailed!} type="error" />}
        </S.Header>
      </Popover>

      <Dropdown overlay={menu}>
        <a onClick={e => e.preventDefault()}>
          <MoreOutlined />
        </a>
      </Dropdown>
    </S.Container>
  );
};

export default AttributeRow;
