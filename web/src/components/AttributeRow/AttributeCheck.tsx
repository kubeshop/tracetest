import {Dropdown, Menu} from 'antd';

import {useRunLayout} from 'components/RunLayout';
import useScrollTo from 'hooks/useScrollTo';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import TraceAnalyticsService from 'services/Analytics/TraceAnalytics.service';
import {IResult} from 'types/Assertion.types';
import * as S from './AttributeRow.styled';

interface IProps {
  items: IResult[];
  type: 'error' | 'success';
}

const AttributeCheck = ({items, type}: IProps) => {
  const {openBottomPanel} = useRunLayout();
  const {setSelectedAssertion} = useTestDefinition();
  const scrollTo = useScrollTo();

  const handleOnClick = (id: string) => {
    TraceAnalyticsService.onAttributeCheckClick();
    const {assertionResult} = items.find(item => item.id === id)!;
    openBottomPanel();
    setSelectedAssertion(assertionResult);
    scrollTo({elementId: `assertion-${id}`, containerId: 'assertions-container'});
  };

  const menuLayout = (
    <Menu
      items={items.map(item => ({
        key: item.id,
        label: item.label,
      }))}
      onClick={({key}) => handleOnClick(key)}
    />
  );

  if (items.length === 1) {
    return (
      <div onClick={() => handleOnClick(items[0].id)}>
        <S.CustomBadge status={type} text={items.length} />
      </div>
    );
  }

  return (
    <div>
      <Dropdown overlay={menuLayout} placement="bottomLeft" trigger={['click']}>
        <div>
          <S.CustomBadge status={type} text={items.length} />
        </div>
      </Dropdown>
    </div>
  );
};

export default AttributeCheck;
