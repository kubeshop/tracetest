import {CaretDownOutlined} from '@ant-design/icons';
import {Menu} from 'antd';
import {useCallback, useMemo} from 'react';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import * as S from './Home.styled';

const {onCreateTestClick} = HomeAnalyticsService;

interface IProps {
  onCreateTest(): void;
  onCreateTransaction(): void;
}

const HomeActions = ({onCreateTest, onCreateTransaction}: IProps) => {
  const onClick = useCallback(
    (key: string) => {
      onCreateTestClick();
      if (key === 'test') return onCreateTest();

      onCreateTransaction();
    },
    [onCreateTest, onCreateTransaction]
  );

  const createMenu = useMemo(
    () => (
      <Menu
        onClick={({key}) => onClick(key)}
        items={[
          {
            label: 'Create New Test',
            key: 'test',
          },
          {
            label: 'Create New Transaction',
            key: 'transaction',
          },
        ]}
      />
    ),
    [onClick]
  );

  return (
    <S.ActionContainer>
      <S.CreateDropdownButton
        overlay={createMenu}
        overlayClassName="test-create-selector-items"
        trigger={['click']}
        placement="bottomRight"
      >
        <S.CreateTestButton type="primary" data-cy="create-button">
          Create <CaretDownOutlined />
        </S.CreateTestButton>
      </S.CreateDropdownButton>
    </S.ActionContainer>
  );
};

export default HomeActions;
