import {TextProps} from 'antd/lib/typography/Text';
import {useMemo, useState} from 'react';

import {isJson} from 'utils/Common';
import Highlighted from '../Highlighted';
import * as S from './AttributeValue.styled';

interface IProps extends TextProps {
  value?: string;
  searchText?: string;
}

const AttributeValue = ({value = '', searchText = '', ...props}: IProps) => {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const isJsonValue = useMemo(() => isJson(value), [value]);

  if (isJsonValue) {
    return (
      <S.ValueJson $isCollapsed={isCollapsed} onClick={() => setIsCollapsed(!isCollapsed)} {...props}>
        <pre>
          <Highlighted highlight={searchText} text={JSON.stringify(JSON.parse(value), null, 2)} />
        </pre>
      </S.ValueJson>
    );
  }

  return (
    <S.ValueText {...props}>
      <Highlighted text={value} highlight={searchText} />
    </S.ValueText>
  );
};

export default AttributeValue;
