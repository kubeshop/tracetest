import {Dropdown, Menu, MenuProps, Space} from 'antd';

import useCopy from 'hooks/useCopy';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import {TSpanFlatAttribute} from 'types/Span.types';
import * as S from './AttributeActions.styled';

const Action = {
  Copy: '0',
  Create_Spec: '1',
  Create_Output: '2',
} as const;

const menuItems: MenuProps['items'] = [
  {
    key: Action.Copy,
    label: 'Copy value',
  },
  {
    key: Action.Create_Output,
    label: 'Create test output',
  },
  {
    key: Action.Create_Spec,
    label: 'Create test spec',
  },
];

interface IProps {
  attribute: TSpanFlatAttribute;
  children?: React.ReactNode;
  onCreateTestOutput(attribute: TSpanFlatAttribute): void;
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
}

const AttributeActions = ({attribute, children, onCreateTestOutput, onCreateTestSpec}: IProps) => {
  const copy = useCopy();

  const handleOnClick = ({key: option}: {key: string}) => {
    if (option === Action.Copy) {
      TraceAnalyticsService.onAttributeCopy();
      copy(attribute.value);
    }

    if (option === Action.Create_Spec) {
      return onCreateTestSpec(attribute);
    }

    if (option === Action.Create_Output) {
      return onCreateTestOutput(attribute);
    }
  };
  return (
    <Dropdown overlay={<Menu items={menuItems} onClick={handleOnClick} />}>
      {children || (
        <a onClick={e => e.preventDefault()}>
          <Space>
            <S.MoreIcon />
          </Space>
        </a>
      )}
    </Dropdown>
  );
};

export default AttributeActions;
