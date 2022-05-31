import {TextProps} from 'antd/lib/typography/Text';
import {useMemo, useState} from 'react';
import JSONPretty from 'react-json-pretty';

import {isJson} from 'utils/Common';
import * as S from './AttributeValue.styled';

interface IProps extends TextProps {
  value: string;
}

const AttributeValue: React.FC<IProps> = ({value, ...props}) => {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const isJsonValue = useMemo(() => isJson(value), [value]);

  if (isJsonValue) {
    return (
      <S.ValueJson isCollapsed={isCollapsed} onClick={() => setIsCollapsed(!isCollapsed)} {...props}>
        <JSONPretty data={value} />
      </S.ValueJson>
    );
  }

  return <S.ValueText {...props}>{value}</S.ValueText>;
};

export default AttributeValue;
