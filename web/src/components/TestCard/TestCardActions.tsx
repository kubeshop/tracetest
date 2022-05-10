import {Dropdown, Menu} from 'antd';
import * as S from './TestCard.styled';

interface ITestCardActions {
  testId: string;
  onDelete(): void;
}

const TestCardActions: React.FC<ITestCardActions> = ({testId, onDelete}) => {
  return (
    <Dropdown
      overlay={
        <Menu>
          <Menu.Item
            data-cy="test-delete-button"
            onClick={({domEvent}) => {
              domEvent.stopPropagation();
              onDelete();
            }}
            key="delete"
          >
            Delete
          </Menu.Item>
        </Menu>
      }
      placement="bottomLeft"
      trigger={['click']}
    >
      <span data-cy={`test-actions-button-${testId}`} className="ant-dropdown-link" onClick={e => e.stopPropagation()}>
        <S.ActionButton />
      </span>
    </Dropdown>
  );
};

export default TestCardActions;
