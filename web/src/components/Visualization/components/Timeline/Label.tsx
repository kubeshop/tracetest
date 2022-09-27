import {Group} from '@visx/group';

import settingIcon from 'assets/setting.svg';
import {SemanticGroupNames, SemanticGroupNamesToText} from 'constants/SemanticGroupNames.constants';
import {SpanKind, SpanKindToText} from 'constants/Span.constants';
import * as S from './Timeline.styled';

interface IProps {
  duration: string;
  kind: SpanKind;
  name: string;
  service: string;
  system: string;
  totalFailedChecks?: number;
  totalPassedChecks?: number;
  type: SemanticGroupNames;
}

const Label = ({duration, kind, name, service, system, totalFailedChecks, totalPassedChecks, type}: IProps) => {
  const failedChecksX = totalPassedChecks ? 20 : 0;

  return (
    <Group left={32} top={0}>
      <Group left={25} top={0}>
        <S.RectBadge rx={2} x={-25} y={0} $type={type} />
        <S.TextBadge dominantBaseline="hanging" textAnchor="middle" y={3}>
          {SemanticGroupNamesToText[type]}
        </S.TextBadge>
      </Group>

      <Group left={60} top={6}>
        {Boolean(totalPassedChecks) && (
          <>
            <S.CircleCheck cx={0} cy={0} r={4} $passed />
            <S.TextDescription dominantBaseline="middle" x={6} y={0}>
              {totalPassedChecks}
            </S.TextDescription>
          </>
        )}
        {Boolean(totalFailedChecks) && (
          <>
            <S.CircleCheck cx={failedChecksX} cy={0} r={4} $passed={false} />
            <S.TextDescription dominantBaseline="middle" x={failedChecksX + 6} y={0}>
              {totalFailedChecks}
            </S.TextDescription>
          </>
        )}
      </Group>

      <Group left={0} top={16}>
        <S.TextName dominantBaseline="hanging" x={0} y={0}>
          {name}
        </S.TextName>

        <S.Image href={settingIcon} y={14} />

        <S.TextDescription dominantBaseline="hanging" dy="1.3em" x={12}>
          <tspan>{`${service} ${SpanKindToText[kind]}`}</tspan>
          {Boolean(system) && <tspan>{` - ${system}`}</tspan>}
          <tspan>{` - ${duration}`}</tspan>
        </S.TextDescription>
      </Group>
    </Group>
  );
};

export default Label;
