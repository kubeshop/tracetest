import {Group} from '@visx/group';

import settingIcon from 'assets/setting.svg';
import {SemanticGroupNames, SemanticGroupNamesToText} from 'constants/SemanticGroupNames.constants';
import {SpanKind, SpanKindToText} from 'constants/Span.constants';
import * as S from '../Timeline.styled';

interface IProps {
  duration: string;
  header?: React.ReactNode;
  kind: SpanKind;
  name: string;
  service: string;
  system: string;
  type: SemanticGroupNames;
}

const Label = ({duration, header, kind, name, service, system, type}: IProps) => (
  <Group left={32} top={0}>
    <Group left={25} top={0}>
      <S.RectBadge rx={2} x={-25} y={0} $type={type} />
      <S.TextBadge dominantBaseline="hanging" textAnchor="middle" y={3}>
        {SemanticGroupNamesToText[type]}
      </S.TextBadge>
    </Group>

    <Group left={60} top={6}>
      {header}
    </Group>

    <Group left={0} top={16}>
      <S.TextName dominantBaseline="hanging" x={0} y={0}>
        {name}
      </S.TextName>

      <S.Image href={settingIcon} y={14} />

      <S.TextDescription dominantBaseline="hanging" dy="1.3em" x={12}>
        <tspan dominantBaseline="hanging">{`${service} ${SpanKindToText[kind]}`}</tspan>
        {Boolean(system) && <tspan dominantBaseline="hanging">{` - ${system}`}</tspan>}
        <tspan dominantBaseline="hanging">{` - ${duration}`}</tspan>
      </S.TextDescription>
    </Group>
  </Group>
);

export default Label;
