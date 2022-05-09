import {Typography} from 'antd';
import SkeletonTable from 'components/SkeletonTable';
import {FC} from 'react';
import {ISpanFlatAttribute} from 'types/Span.types';
import * as S from './Attributes.styled';
import SpanAttributesTable from '../../SpanAttributesTable/SpanAttributesTable';

interface TAttributesProps {
  spanAttributeList?: ISpanFlatAttribute[];
}

export const Attributes: FC<TAttributesProps> = ({spanAttributeList = []}) => {
  return (
    <S.Container data-cy="span-details-attributes">
      <S.Header>
        <Typography.Text strong>Attributes</Typography.Text>
      </S.Header>
      <SkeletonTable loading={!spanAttributeList.length}>
        <SpanAttributesTable spanAttributesList={spanAttributeList} />
      </SkeletonTable>
    </S.Container>
  );
};
