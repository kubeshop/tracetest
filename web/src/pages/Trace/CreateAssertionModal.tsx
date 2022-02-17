import {useState} from 'react';
import styled from 'styled-components';
import {Button, Input, List, Modal, Select as AntSelect, AutoComplete, Typography} from 'antd';

import {ISpan, ISpanTag} from '../../types';

const {Title} = Typography;
interface IProps {
  open: boolean;
  onClose: () => void;
  span: ISpan;
}

const Select = styled(AntSelect)`
  min-width: 88px;
  > .ant-select-selector {
    min-height: 100%;
  }
`;

const CreateAssertionModal = ({span, open, onClose}: IProps) => {
  const [assertionList, setAssertionList] = useState<any>(Array(3).fill({}));

  const spanTagsMap = span.tags.reduce<{[key: string]: ISpanTag[]}>((acc, item) => {
    const keyPrefix = item.key.split('.').shift();
    if (!keyPrefix) {
      return acc;
    }
    const keys = acc[keyPrefix] || [];
    keys.push(item);
    acc[keyPrefix] = keys;
    return acc;
  }, {});

  const handleAddItem = () => {
    setAssertionList([...assertionList, {key: '', value: '', type: ''}]);
  };

  const renderTitle = (title: string) => <span>{title}</span>;

  const renderItem = (title: string) => ({
    value: title,
    label: (
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
        }}
      >
        {title}
      </div>
    ),
  });

  const keysOptions = Object.keys(spanTagsMap)
    .sort((el1, el2) => el1.charCodeAt(0) - el2.charCodeAt(0))
    .map(tagKey => {
      return {
        label: renderTitle(tagKey),
        options: spanTagsMap[tagKey].map(el => renderItem(el.key)),
      };
    });

  return (
    <Modal style={{minHeight: 500}} footer={null} visible={span && open} onCancel={onClose}>
      <Title level={2}>Create New Assertion </Title>
      <List
        dataSource={assertionList}
        itemLayout="horizontal"
        loadMore={
          <Button style={{marginTop: 8}} onClick={handleAddItem}>
            Add Item
          </Button>
        }
        renderItem={(item: any, index) => {
          return (
            <div key={`key_${index}`} style={{display: 'flex', alignItems: 'stretch', gap: 4, marginBottom: 4}}>
              <AutoComplete style={{width: 250}} options={keysOptions}>
                <Input.Search size="large" placeholder="span key" />
              </AutoComplete>
              <Select defaultValue="Op1">
                <Select.Option value="eq">eq</Select.Option>
                <Select.Option value="ne">ne</Select.Option>
                <Select.Option value="gt">gt</Select.Option>
                <Select.Option value="lt">lt</Select.Option>
                <Select.Option value="ge">ge</Select.Option>
                <Select.Option value="le">le</Select.Option>
              </Select>
              <Input size="small" style={{width: 'unset'}} placeholder="type value" />
            </div>
          );
        }}
      />
      <div style={{padding: 16, display: 'flex', justifyContent: 'flex-end'}}>
        <Button>Create</Button>
      </div>
    </Modal>
  );
};

export default CreateAssertionModal;
