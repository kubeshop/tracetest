import {CaretDownOutlined} from '@ant-design/icons';
import {Menu} from 'antd';
import {useCallback, useMemo} from 'react';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import * as S from './Home.styled';

const {onCreateTestClick} = HomeAnalyticsService;

interface IProps {
  onCreateTest(): void;
  onCreateTestSuite(): void;
}

const HomeActions = ({onCreateTest, onCreateTestSuite}: IProps) => {
  const onClick = useCallback(
    (key: string) => {
      onCreateTestClick();
      if (key === 'test') return onCreateTest();

      onCreateTestSuite();
    },
    [onCreateTest, onCreateTestSuite]
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
            label: 'Create New Test Suite',
            key: 'testsuite',
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
