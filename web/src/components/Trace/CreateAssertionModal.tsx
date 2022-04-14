import {ChangeEvent, useState} from 'react';
import {QuestionCircleOutlined, PlusOutlined, TagOutlined} from '@ant-design/icons';
import {isEmpty} from 'lodash';
import styled from 'styled-components';
import {Button, Input, List, Modal, Select as AntSelect, AutoComplete, Typography, Tooltip, Tag, Checkbox} from 'antd';
import jemsPath from 'jmespath';

import {COMPARE_OPERATOR, ISpan, ItemSelector, ITrace, LOCATION_NAME, SpanSelector} from 'types';
import {useCreateAssertionMutation} from 'services/TestService';
import {SELECTOR_DEFAULT_ATTRIBUTES} from 'lib/SelectorDefaultAttributes';
import {filterBySpanId} from 'utils';
import {getSpanSignature} from 'services/SpanService';
import {getEffectedSpansCount} from 'services/AssertionService';

interface IItemSelectorDropdown {
  span: ISpan;
  trace: ITrace;
  itemSelectors: Array<ItemSelector>;
  onChangeItemSelector: (selectors: Array<ItemSelector>) => void;
}

interface ISpanSelectorOption {
  key: string;
  compareOp: keyof typeof COMPARE_OPERATOR;
  value: string;
}

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

const ItemSelectorDropdown = ({span, trace, itemSelectors, onChangeItemSelector}: IItemSelectorDropdown) => {
  const [itemSelectorInput, setItemSelectorInput] = useState('');

  const itemSelectorsList = getSpanSignature(span.spanId, trace);
  const itemSelectorOptions = itemSelectorsList.map((tag: any, index) => {
    return {
      label: (
        <span>
          <Checkbox
            style={{marginLeft: 8, marginRight: 8}}
            checked={itemSelectors.findIndex(el => el.propertyName.includes(tag.propertyName)) > -1}
          />{' '}
          {`${tag.propertyName} (${tag.value})`}
        </span>
      ),
      value: tag.propertyName,
    };
  });

  const handleSelectItemSelector = (text: any) => {
    if (itemSelectors.findIndex(el => el.propertyName.includes(text)) > -1) {
      onChangeItemSelector(itemSelectors.filter(el => !el.propertyName.includes(text)));
    } else {
      const selectedItem = itemSelectorsList.find(el => el.propertyName.includes(text));
      if (selectedItem) {
        onChangeItemSelector([...itemSelectors, selectedItem]);
      }
    }
    setItemSelectorInput('');
  };

  const handleDeleteItemSelector = (item: ItemSelector) => {
    onChangeItemSelector(itemSelectors.filter(el => el.propertyName !== item.propertyName));
  };

  return (
    <>
      <div style={{marginBottom: 8, display: 'flex', flexWrap: 'wrap'}}>
        {itemSelectors.map((item: ItemSelector) => (
          <Tag key={item.propertyName} closable onClose={() => handleDeleteItemSelector(item)}>
            {item.value}
          </Tag>
        ))}
      </div>

      <AutoComplete
        style={{width: '100%'}}
        options={itemSelectorOptions}
        onSelect={handleSelectItemSelector}
        searchValue={itemSelectorInput}
        value={itemSelectorInput}
        onSearch={setItemSelectorInput}
        backfill
        filterOption={(inputValue, option) => {
          return Boolean(option?.value.includes(inputValue));
        }}
      >
        <Input prefix={<TagOutlined style={{marginRight: 4}} />} size="large" placeholder="Add a selector" />
      </AutoComplete>
    </>
  );
};

const itemSelectorKeys = SELECTOR_DEFAULT_ATTRIBUTES.map(el => el.attributes).flat();

const initialFormState = Array(1).fill({key: '', compareOp: COMPARE_OPERATOR.EQUALS, value: ''});

const effectedSpanMessage = (spanCount: number) => {
  if (spanCount <= 1) {
    return `Effects ${spanCount} span`;
  }

  return `Effects ${spanCount} spans`;
};

const CreateAssertionModal = ({testId, span, trace, open, onClose}: IProps) => {
  const [assertionList, setAssertionList] = useState<Array<Partial<ISpanSelectorOption>>>(initialFormState);
  const [itemSelectors, setItemSelectors] = useState<ItemSelector[]>([]);
  const [createAssertion] = useCreateAssertionMutation();
  const attrs = jemsPath.search(trace, filterBySpanId(span.spanId));
  const effectedSpanCount = getEffectedSpansCount(trace, itemSelectors);
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

  const spanAssertionOptions = Object.keys(spanTagsMap).map((tagKey: any, index) => {
    return {
      label: renderTitle(tagKey, index),
      options: spanTagsMap[tagKey].map((el: any) => renderItem(el)),
    };
  });

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
      title={
        <div style={{display: 'flex', justifyContent: 'space-between', marginRight: 36}}>
          <Typography.Title level={5}>Create New Assertion</Typography.Title>
          <Typography.Text>{effectedSpanMessage(effectedSpanCount)}</Typography.Text>
        </div>
      }
      onOk={handleCreateAssertion}
      okButtonProps={{
        type: 'default',
        disabled: !isValid,
      }}
      okText="Save"
    >
      <div style={{marginBottom: 8}}>
        <Typography.Text style={{marginRight: 8}}>Selectors</Typography.Text>
        <Tooltip title="prompt text">
          <QuestionCircleOutlined style={{color: '#8C8C8C'}} />
        </Tooltip>
      </div>

      <ItemSelectorDropdown
        span={span}
        trace={trace}
        itemSelectors={itemSelectors}
        onChangeItemSelector={setItemSelectors}
      />
      <div style={{marginTop: 24, marginBottom: 8}}>
        <Typography.Text style={{marginRight: 8}}>Span Assertions</Typography.Text>
        <Tooltip title="prompt text">
          <QuestionCircleOutlined style={{color: '#8C8C8C'}} />
        </Tooltip>
      </div>
      <List
        dataSource={assertionList}
        itemLayout="horizontal"
        loadMore={
          <Button type="link" icon={<PlusOutlined />} style={{padding: 0}} onClick={handleAddItem}>
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
                options={spanAssertionOptions}
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
