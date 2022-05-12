import {Dropdown, Menu} from 'antd';
import * as S from './ResultCard.styled';

interface IResultCardActionsProps {
  resultId: string;
  onDelete(resultId: string): void;
}

const ResultCardActions: React.FC<IResultCardActionsProps> = ({resultId, onDelete}) => {
  return (
    <span
      data-cy={`result-actions-button-${resultId}`}
      className="ant-dropdown-link"
      onClick={e => e.stopPropagation()}
      style={{textAlign: 'right'}}
    >
      <Dropdown
        overlay={
          <Menu>
            <Menu.Item
              data-cy="test-delete-button"
              onClick={({domEvent}) => {
                domEvent.stopPropagation();
                onDelete(resultId);
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
        <S.ActionButton />
      </Dropdown>
    </span>
  );
};

export default ResultCardActions;
