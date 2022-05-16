import {TextProps} from 'antd/lib/typography/Text';
import {useMemo, useState} from 'react';
import JSONPretty from 'react-json-pretty';
import {isJson} from '../../utils/Common';
import * as S from './AttributeValue.styled';

interface IAttributeValueProps extends TextProps {
  value: string;
}

const AttributeValue: React.FC<IAttributeValueProps> = ({value, ...props}) => {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const parsedValue = useMemo(() => (isJson(value) ? <JSONPretty data={value} /> : value), [value]);

  return (
    <S.ValueText onClick={() => setIsCollapsed(!isCollapsed)} isCollapsed={isCollapsed} {...props}>
      {parsedValue}
    </S.ValueText>
  );
};

export default AttributeValue;
