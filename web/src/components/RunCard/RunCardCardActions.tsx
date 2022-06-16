import {Dropdown, Menu} from 'antd';
import * as S from './RunCard.styled';

interface IProps {
  resultId: string;
  testId: string;
  onDelete(resultId: string): void;
}

const ResultCardActions = ({resultId, testId, onDelete}: IProps) => {
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
            <Menu.Item data-cy="download-junit-button" key="download-junit">
              <a href={`/api/tests/${testId}/run/${resultId}/junit.xml`} download>
                Download JUnit
              </a>
            </Menu.Item>
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
