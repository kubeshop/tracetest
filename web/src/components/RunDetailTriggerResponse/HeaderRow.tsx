import {MoreOutlined} from '@ant-design/icons';
import {Dropdown, Menu, message} from 'antd';
import {THeader} from 'types/Test.types';
import AttributeValue from '../AttributeValue';
import * as S from './RunDetailTriggerResponse.styled';

interface IProps {
  header: THeader;
  onCopy(value: string): void;
}

const HeaderRow = ({header: {key = '', value = ''}, onCopy}: IProps) => {
  const handleOnClick = () => {
    message.success('Value copied to the clipboard');
    return onCopy(value);
  };

  const menu = (
    <Menu
      items={[
        {
          label: 'Copy value',
          key: 'copy',
        },
      ]}
      onClick={handleOnClick}
    />
  );

  return (
    <S.HeaderContainer>
      <S.Header>
        <S.AttributeTitle title={key} />
        <S.AttributeValueRow>
          <AttributeValue value={value} />
        </S.AttributeValueRow>
      </S.Header>

      <Dropdown overlay={menu} trigger={['click']}>
        <a onClick={e => e.preventDefault()}>
          <MoreOutlined />
        </a>
      </Dropdown>
    </S.HeaderContainer>
  );
};

export default HeaderRow;
