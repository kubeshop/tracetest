import {ChangeEvent, useState} from 'react';
import {isEmpty} from 'lodash';
import styled from 'styled-components';
import {Button, Input, List, Modal, Select as AntSelect, AutoComplete, Typography} from 'antd';
import jemsPath from 'jmespath';
import Text from 'antd/lib/typography/Text';

import {COMPARE_OPERATOR, IAttribute, ISpan, LOCATION_NAME, SpanSelector} from 'types';
import {useCreateAssertionMutation} from 'services/TestService';
import {SELECTOR_DEFAULT_ATTRIBUTES} from 'lib/SelectorDefaultAttributes';
import {filterBySpanId} from 'utils';
import {getSpanSignature} from '../../services/SpanService';

const {Title} = Typography;
interface IProps {
  open: boolean;
  onClose: () => void;
  span: ISpan;
  testId: string;
  trace: any;
}

const Select = styled(AntSelect)`
  min-width: 88px;
  > .ant-select-selector {
    min-height: 100%;
  }
`;

const NORMALIZED_OBJECT = `resourceSpans[].{instrumentationLibrarySpans:instrumentationLibrarySpans[].{spans:spans[].{spanId:spanId,attributes:attributes[].{key:key,value:value.*|[0],type:'span'}}},resource:resource.{attributes:attributes[].{key:key,value: value.*|[0],type:'resource'}}}`;
const filterByAttributes = (condition: string) => `[${condition && '?'}${condition}]`;

const selectionPipe = (query: string) => `${NORMALIZED_OBJECT} | ${query}  | length([])`;

const createSelector = (query: string) => `${NORMALIZED_OBJECT} | ${query}`;

const selectorConditionBuilder = (attribute: IAttribute) => {
  if (attribute.type === 'span') {
    return `instrumentationLibrarySpans[?spans[?attributes[?key=='${attribute.key}' && value=='${attribute.value}']]]`;
  }

  if (attribute.type === 'resource') {
    return `resource.attributes[?key=='${attribute.key}' && value=='${attribute.value}']`;
  }
};

const itemSelectorKeys = SELECTOR_DEFAULT_ATTRIBUTES.map(el => el.attributes).flat();

const initialFormState = Array(3).fill({key: '', compareOp: COMPARE_OPERATOR.EQUALS, value: ''});

const CreateAssertionModal = ({testId, span, trace, open, onClose}: IProps) => {
  const [assertionList, setAssertionList] =
    useState<Array<Partial<{key: string; compareOp: keyof typeof COMPARE_OPERATOR; value: string}>>>(initialFormState);
  const [createAssertion] = useCreateAssertionMutation();
  const attrs = jemsPath.search(trace, filterBySpanId(span.spanId));

  const selectorCondition = assertionList
    .map(k => {
      return attrs?.find((el: any) => el.key === k.key);
    })
    .filter(i => i)
    .map(item => selectorConditionBuilder(item))
    .join(' && ');

  const effectedSpans =
    selectorCondition.length > 0 ? jemsPath.search(trace, selectionPipe(filterByAttributes(selectorCondition))) : 0;

  const spanTagsMap = attrs?.reduce((acc: {[x: string]: any}, item: {key: string}) => {
    if (itemSelectorKeys.indexOf(item.key) !== -1) {
      return acc;
    }
    const keyPrefix = item.key.split('.').shift() || item.key;
    if (!keyPrefix) {
      return acc;
    }
    const keys = acc[keyPrefix] || [];
    keys.push(item);
    acc[keyPrefix] = keys;
    return acc;
  }, {});

  const handleAddItem = () => {
    setAssertionList([...assertionList, {key: '', compareOp: COMPARE_OPERATOR.EQUALS}]);
  };

  const renderTitle = (title: any, index: number) => (
    <span key={`KEY_${title}_${index}`}>{`${title} (${spanTagsMap[title][0].type})`}</span>
  );

  const renderItem = (attr: any) => ({
    value: attr.key,
    label: (
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
        }}
      >
        {attr.key}
      </div>
    ),
  });

  const keysOptions = Object.keys(spanTagsMap).map((tagKey: any, index) => {
    return {
      label: renderTitle(tagKey, index),
      options: spanTagsMap[tagKey].map((el: any) => renderItem(el)),
    };
  });

  const itemSelectors = getSpanSignature(span.spanId, trace);

  const handleClose = () => {
    setAssertionList(initialFormState);
    onClose();
  };

  const handleCreateAssertion = async () => {
    const spanAssertions = assertionList
      .map(k => ({value: k.value, attr: attrs?.find((el: any) => el.key === k?.key), compareOp: k.compareOp}))
      .filter(i => i.attr)
      .map(({attr, value, compareOp}) => {
        const spanAttribute = span.attributes.find(({key}) => key.includes(attr.key));

        return {
          locationName: attr.type.includes('span') ? LOCATION_NAME.SPAN_ATTRIBUTES : LOCATION_NAME.RESOURCE_ATTRIBUTES,
          propertyName: attr.key as string,
          comparisonValue: value,
          operator: compareOp,
          valueType: spanAttribute?.value.intValue ? 'intValue' : 'stringValue',
        };
      })
      .filter((el): el is SpanSelector => Boolean(el));

    createAssertion({testId, selectors: itemSelectors, spanAssertions});
    handleClose();
  };

  return (
    <Modal style={{minHeight: 500}} footer={null} visible={span && open} onCancel={handleClose} destroyOnClose>
      <Title level={2}>Create New Assertion</Title>
      <div style={{display: 'flex', justifyContent: 'flex-end', height: 64}}>
        <Text>Effects {effectedSpans} spans</Text>
      </div>
      <List
        dataSource={assertionList}
        itemLayout="horizontal"
        loadMore={
          <Button style={{marginTop: 8}} onClick={handleAddItem}>
            Add Item
          </Button>
        }
        renderItem={(item: any, index) => {
          const handleSearchChange = (searchText: string) => {
            const value = attrs?.find((el: any) => el.key === assertionList[index].key)?.value;

            assertionList[index] = {...assertionList[index], key: searchText, value: isEmpty(value) ? '' : value};
            setAssertionList([...assertionList]);
          };

          const handleSelectCompareOperator = (compareOp: any) => {
            assertionList[index] = {...assertionList[index], compareOp};
            setAssertionList([...assertionList]);
          };

          const handleValueChange = (event: ChangeEvent<HTMLInputElement>) => {
            assertionList[index] = {...assertionList[index], value: event.target.value};
            setAssertionList([...assertionList]);
          };

          return (
            <div key={`key_${index}`} style={{display: 'flex', alignItems: 'stretch', gap: 4, marginBottom: 4}}>
              <AutoComplete
                style={{width: 250}}
                onChange={handleSearchChange}
                onSearch={handleSearchChange}
                onSelect={handleSearchChange}
                options={keysOptions}
                filterOption={(inputValue, option) => {
                  return option?.label.props.children.includes(inputValue);
                }}
              >
                <Input.Search size="large" placeholder="span key" />
              </AutoComplete>
              <Select defaultValue="eq" onSelect={handleSelectCompareOperator}>
                <Select.Option value={COMPARE_OPERATOR.EQUALS}>eq</Select.Option>
                <Select.Option value={COMPARE_OPERATOR.NOTEQUALS}>ne</Select.Option>
                <Select.Option value={COMPARE_OPERATOR.GREATERTHAN}>gt</Select.Option>
                <Select.Option value={COMPARE_OPERATOR.LESSTHAN}>lt</Select.Option>
                <Select.Option value={COMPARE_OPERATOR.GREATOREQUALS}>ge</Select.Option>
                <Select.Option value={COMPARE_OPERATOR.LESSOREQUAL}>le</Select.Option>
              </Select>
              <Input
                size="small"
                style={{width: 'unset'}}
                placeholder="type value"
                onChange={handleValueChange}
                value={assertionList[index].value}
              />
            </div>
          );
        }}
      />
      <div style={{padding: 16, display: 'flex', justifyContent: 'flex-end'}}>
        <Button onClick={handleCreateAssertion}>Create</Button>
      </div>
    </Modal>
  );
};

export default CreateAssertionModal;
