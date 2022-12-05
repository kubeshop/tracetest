import {Dropdown, Menu} from 'antd';
import {uniqBy} from 'lodash';

import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import {IResult} from 'types/Assertion.types';
import * as S from './AssertionResultChecks.styled';

interface IProps {
  items: IResult[];
  type: 'error' | 'success';
  styleType?: 'node' | 'summary' | 'default';
}

const ResultCheck = ({items, type, styleType = 'default'}: IProps) => {
  const {setSelectedSpec} = useTestSpecs();

  const handleOnClick = (id: string) => {
    TraceAnalyticsService.onAttributeCheckClick();
    const {assertionResult} = items.find(item => item.id === id)!;
    setSelectedSpec(assertionResult.selector);
  };

  const menuLayout = (
    <Menu
      items={uniqBy(items, 'id').map(item => ({
        key: item.id,
        label: item.label,
      }))}
      onClick={({key}) => handleOnClick(key)}
    />
  );

  if (items.length === 1) {
    return (
      <div onClick={() => handleOnClick(items[0].id)}>
        <S.CustomBadge status={type} text={items.length} $styleType={styleType} />
      </div>
    );
  }

  return (
    <div>
      <Dropdown overlay={menuLayout} placement="bottomLeft" trigger={['click']}>
        <S.CustomBadge status={type} text={items.length} $styleType={styleType} />
      </Dropdown>
    </div>
  );
};

export default ResultCheck;
