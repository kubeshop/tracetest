import {MoreOutlined} from '@ant-design/icons';
import {Dropdown, Menu, message} from 'antd';

import {IResult} from 'types/Assertion.types';
import {TSpanFlatAttribute} from 'types/Span.types';
import AttributeValue from '../AttributeValue';
import AttributeCheck from './AttributeCheck';
import * as S from './AttributeRow.styled';

interface IProps {
  assertionsFailed?: IResult[];
  assertionsPassed?: IResult[];
  attribute: TSpanFlatAttribute;
  onCopy(value: string): void;
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
  searchText?: string;
}

enum Action {
  Copy = '0',
  Create_Spec = '1',
}

const AttributeRow = ({
  assertionsFailed,
  assertionsPassed,
  attribute: {key, value},
  attribute,
  onCopy,
  onCreateTestSpec,
  searchText,
}: IProps) => {
  const passedCount = assertionsPassed?.length ?? 0;
  const failedCount = assertionsFailed?.length ?? 0;

  const handleOnClick = ({key: option}: {key: string}) => {
    if (option === Action.Copy) {
      message.success('Value copied to the clipboard');
      return onCopy(value);
    }

    if (option === Action.Create_Spec) {
      return onCreateTestSpec(attribute);
    }
  };

  const menu = (
    <Menu
      items={[
        {
          label: 'Copy value',
          key: Action.Copy,
        },
        {
          label: 'Create test spec',
          key: Action.Create_Spec,
        },
      ]}
      onClick={handleOnClick}
    />
  );

  return (
    <S.Container>
      <S.Header>
        <S.AttributeTitle title={key} searchText={searchText} />
        <S.AttributeValueRow>
          <AttributeValue value={value} searchText={searchText} />
        </S.AttributeValueRow>
        {passedCount > 0 && <AttributeCheck items={assertionsPassed!} type="success" />}
        {failedCount > 0 && <AttributeCheck items={assertionsFailed!} type="error" />}
      </S.Header>

      <Dropdown overlay={menu}>
        <a onClick={e => e.preventDefault()}>
          <MoreOutlined />
        </a>
      </Dropdown>
    </S.Container>
  );
};

export default AttributeRow;
