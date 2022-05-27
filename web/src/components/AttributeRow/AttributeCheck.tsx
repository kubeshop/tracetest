import {Dropdown, Menu} from 'antd';

import {IResult} from 'components/SpanDetail/SpanDetail';
import * as S from './AttributeRow.styled';

interface IProps {
  items: IResult[];
  type: 'error' | 'success';
}

const AttributeCheck = ({items, type}: IProps) => {
  const handleOnClick = (id: string) => {
    console.log('### click', id);
    const resultListElement = document.getElementById(`assertion-${id}`);
    const offsetTop = resultListElement?.offsetTop ?? 0;
    console.log('### offsetTop', offsetTop);
    const assertionsContainer = document.getElementById('assertions-container');
    if (assertionsContainer) {
      assertionsContainer.scrollTop = offsetTop;
    }
  };

  const menuLayout = (
    <Menu
      items={items.map(item => ({
        key: item.id,
        label: item.id,
      }))}
      onClick={({key: id}) => handleOnClick(id)}
    />
  );

  if (items.length === 1) {
    return (
      <div onClick={() => handleOnClick(items[0].id)}>
        <S.CustomBadge status={type} text={items.length} />
      </div>
    );
  }

  return (
    <div>
      <Dropdown overlay={menuLayout} placement="bottomLeft" trigger={['click']}>
        <div>
          <S.CustomBadge status={type} text={items.length} />
        </div>
      </Dropdown>
    </div>
  );
};

export default AttributeCheck;
