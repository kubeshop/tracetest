import {Dropdown, Menu} from 'antd';
import {uniqBy} from 'lodash';

import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import {TTestSpec} from 'types/TestRun.types';
import * as S from './AssertionResultChecks.styled';

interface IProps {
  items: TTestSpec[];
  type: 'error' | 'success';
  styleType?: 'node' | 'summary' | 'default';
}

const ResultCheck = ({items, type, styleType = 'default'}: IProps) => {
  const {setSelectedSpec} = useTestSpecs();

  const handleOnClick = (selector: string) => {
    TraceAnalyticsService.onAttributeCheckClick();
    setSelectedSpec(selector);
  };

  const menuLayout = (
    <Menu
      items={uniqBy(items, 'selector').map(item => ({
        key: item.selector,
        label: item.selector,
      }))}
      onClick={({key}) => handleOnClick(key)}
    />
  );

  if (items.length === 1) {
    return (
      <div onClick={() => handleOnClick(items[0].selector)}>
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
