import {Typography} from 'antd';
import SkeletonTable from 'components/SkeletonTable';
import {FC} from 'react';
import {ISpanFlatAttribute} from '../../types/Span.types';
import SpanAttributesTable from '../SpanAttributesTable/SpanAttributesTable';
import * as S from './Attributes.styled';

type TAttributesProps = {
  spanAttributeList?: ISpanFlatAttribute[];
};

const Attributes: FC<TAttributesProps> = ({spanAttributeList = []}) => {
  return (
    <S.Container>
      <S.Header>
        <Typography.Text strong>Attributes</Typography.Text>
      </S.Header>
      <SkeletonTable loading={!spanAttributeList.length}>
        <SpanAttributesTable spanAttributesList={spanAttributeList} />
      </SkeletonTable>
    </S.Container>
  );
};

export default Attributes;
