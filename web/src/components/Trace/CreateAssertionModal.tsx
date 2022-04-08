import {ChangeEvent, useState} from 'react';
import {PlusOutlined} from '@ant-design/icons';
import {isEmpty} from 'lodash';
import styled from 'styled-components';
import {Button, Input, List, Modal, Select as AntSelect, AutoComplete, Typography} from 'antd';
import jemsPath from 'jmespath';

import {COMPARE_OPERATOR, ISpan, LOCATION_NAME, SpanSelector} from 'types';
import {useCreateAssertionMutation} from 'services/TestService';
import {SELECTOR_DEFAULT_ATTRIBUTES} from 'lib/SelectorDefaultAttributes';
import {filterBySpanId} from 'utils';
import {getSpanSignature} from '../../services/SpanService';

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

const itemSelectorKeys = SELECTOR_DEFAULT_ATTRIBUTES.map(el => el.attributes).flat();

const initialFormState = Array(3).fill({key: '', compareOp: COMPARE_OPERATOR.EQUALS, value: ''});

const CreateAssertionModal = ({testId, span, trace, open, onClose}: IProps) => {
  const [assertionList, setAssertionList] =
    useState<Array<Partial<{key: string; compareOp: keyof typeof COMPARE_OPERATOR; value: string}>>>(initialFormState);
  const [createAssertion] = useCreateAssertionMutation();
  const attrs = jemsPath.search(trace, filterBySpanId(span.spanId));

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

  const handleCreateAssertion = async () => {
    createAssertion({testId, selectors: itemSelectors, spanAssertions});
    handleClose();
  };

  const isValid = spanAssertions.length && spanAssertions.length;

  return (
    <Modal
      style={{minHeight: 500}}
      visible={span && open}
      onCancel={handleClose}
      destroyOnClose
      title={<Typography.Title level={5}>Create New Assertion</Typography.Title>}
      onOk={handleCreateAssertion}
      okButtonProps={{
        type: 'default',
        disabled: !isValid,
      }}
      okText="Create"
    >
      <List
        dataSource={assertionList}
        itemLayout="horizontal"
        loadMore={
          <Button type="link" icon={<PlusOutlined />} style={{marginTop: 8, padding: 0}} onClick={handleAddItem}>
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
            <div key={`key_${index}`} style={{display: 'flex', alignItems: 'stretch', gap: 4, marginBottom: 16}}>
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
    </Modal>
  );
};

export default CreateAssertionModal;
