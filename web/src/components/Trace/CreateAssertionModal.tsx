import {useState} from 'react';
import styled from 'styled-components';
import {Button, Input, List, Modal, Select as AntSelect, AutoComplete, Typography} from 'antd';
import jemsPath from 'jmespath';
import Text from 'antd/lib/typography/Text';

import {IAttribute, ISpan} from 'types';
import {useCreateAssertionMutation} from 'services/TestService';

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

const flattenAttributesSelector = () =>
  `resourceSpans[].[instrumentationLibrarySpans[].spans[].attributes[].{key:key,value:value.*|[0]},resource.attributes[].{key:key,value: value.*|[0]}]|[][]`;

const filterBySpanId = (spanId: string = '') =>
  `resourceSpans[?instrumentationLibrarySpans[?spans[?starts_with(spanId,'${spanId}')]]] | [].[instrumentationLibrarySpans[].spans[].attributes[].{key:key,value:value.*|[0],type:'span'},resource.attributes[].{key:key,value: value.*|[0],type:'resource'}]|[][]`;
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

const CreateAssertionModal = ({testId, span, trace, open, onClose}: IProps) => {
  const [assertionList, setAssertionList] = useState<Array<string>>(Array(3).fill(''));
  const [createAssertion, result] = useCreateAssertionMutation();
  const attrs = jemsPath.search(trace, filterBySpanId(span.spanId));
  console.log('@@CreateAssertionModal', attrs);

  const selectorCondition = assertionList
    .map(k => {
      return attrs?.find((el: any) => el.key === k);
    })
    .filter(i => i)
    .map(item => selectorConditionBuilder(item))
    .join(' && ');
  console.log('@@query', selectionPipe(filterByAttributes(selectorCondition)));
  const effectedSpans = jemsPath.search(trace, selectionPipe(filterByAttributes(selectorCondition)));
  const spanTagsMap = attrs.reduce((acc: {[x: string]: any}, item: {key: string}) => {
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
    setAssertionList([...assertionList, '']);
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

  const handleCreateAssertion = async () => {
    await createAssertion({testId, selector: createSelector(filterByAttributes(selectorCondition))});
    onClose();
  };

  return (
    <Modal style={{minHeight: 500}} footer={null} visible={span && open} onCancel={onClose}>
      <Title level={2}>Create New Assertion </Title>
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
            assertionList[index] = searchText;
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
              <Select defaultValue="eq">
                <Select.Option value="eq">eq</Select.Option>
                <Select.Option value="ne">ne</Select.Option>
                <Select.Option value="gt">gt</Select.Option>
                <Select.Option value="lt">lt</Select.Option>
                <Select.Option value="ge">ge</Select.Option>
                <Select.Option value="le">le</Select.Option>
              </Select>
              <Input
                size="small"
                style={{width: 'unset'}}
                placeholder="type value"
                value={attrs?.find((el: any) => el.key === assertionList[index])?.value}
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
