import {useMemo} from 'react';

import {isJson} from 'utils/Common';
import * as S from './OutputRow.styled';

interface IProps {
  value: string;
}

const OutputValue = ({value}: IProps) => {
  const isJsonValue = useMemo(() => isJson(value), [value]);

  if (isJsonValue) {
    return (
      <S.ValueJson>
        <pre>{JSON.stringify(JSON.parse(value), null, 2)}</pre>
      </S.ValueJson>
    );
  }

  return <S.ValueText>{value}</S.ValueText>;
};

export default OutputValue;
