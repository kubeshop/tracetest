import {Dropdown, Menu} from 'antd';
import * as S from './TestHeader.styled';

interface IProps {
  testId: string;
  resultId: string;
}

const Actions = ({testId, resultId}: IProps) => {
  return (
    <span
      data-cy="test-header-actions-button"
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

export default Actions;
