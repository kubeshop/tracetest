import {Typography} from 'antd';
import {FC, useMemo} from 'react';
import {getResourceSpanBySpanId, getSpanAttributeList} from '../../services/SpanService';
import {ITrace} from '../../types';
import SpanAttributesTable from '../SpanAttributesTable/SpanAttributesTable';
import * as S from './Attributes.styled';

type TAttributesProps = {
  spanId?: string;
  trace?: ITrace;
};

const Attributes: FC<TAttributesProps> = ({spanId, trace}) => {
  const spanAttributesList = useMemo(() => {
    if (!spanId || !trace) {
      return [];
    }
    const resourceSpan = getResourceSpanBySpanId(spanId, trace);

    return resourceSpan ? getSpanAttributeList(resourceSpan) : [];
  }, [spanId, trace]);

  return (
    <S.Container>
      <S.Header>
        <Typography.Text strong>Attributes</Typography.Text>
      </S.Header>
      <SpanAttributesTable spanAttributesList={spanAttributesList} />
    </S.Container>
  );
};

export default Attributes;
