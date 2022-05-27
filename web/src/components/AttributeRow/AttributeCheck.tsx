import {Dropdown, Menu} from 'antd';

import {IResult} from 'components/SpanDetail/SpanDetail';
import * as S from './AttributeRow.styled';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';

interface IProps {
  items: IResult[];
  type: 'error' | 'success';
}

const AttributeCheck = ({items, type}: IProps) => {
  const {setIsCollapsed} = useAssertionForm();

  const handleOnClick = (id: string) => {
    setIsCollapsed(true);
    const resultListElement = document.getElementById(`assertion-${id}`);
    const offsetTop = resultListElement?.offsetTop ?? 0;
    const assertionsContainer = document.getElementById('assertions-container');
    if (assertionsContainer) {
      assertionsContainer.scrollTop = offsetTop;
    }
  };

  const menuLayout = (
    <Menu
      items={items.map(item => ({
        key: item.id,
        label: item.label,
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
