import {Dropdown, Menu} from 'antd';
import * as S from './TestCard.styled';
import {useOnDeleteCallback} from './useOnDeleteCallback';

interface IProps {
  testId: string;

  onDelete(): void;
}

const TestCardActions: React.FC<IProps> = ({testId, onDelete}) => {
  const onClick = useOnDeleteCallback(onDelete);
  return (
    <Dropdown
      overlay={<Menu items={[{key: 'delete', label: <p data-cy="delete">Delete</p>, onClick}]} />}
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
