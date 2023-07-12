import {Group} from '@visx/group';
import * as S from '../Timeline.styled';

interface IProps {
  hasOutputs?: boolean;
  totalFailedChecks?: number;
  totalPassedChecks?: number;
}

function getOutputsX(totalFailedChecks?: number, totalPassedChecks?: number): number {
  if (totalFailedChecks && totalPassedChecks) {
    return 44;
  }
  if (totalFailedChecks || totalPassedChecks) {
    return 24;
  }

  return 0;
}

const Header = ({hasOutputs, totalFailedChecks, totalPassedChecks}: IProps) => {
  const failedChecksX = totalPassedChecks ? 20 : 0;
  const outputsX = getOutputsX(totalFailedChecks, totalPassedChecks);

  return (
    <>
      <Group>
        {!!totalPassedChecks && (
          <>
            <S.CircleCheck cx={0} cy={0} r={4} $passed />
            <S.TextDescription dominantBaseline="middle" x={6} y={0}>
              {totalPassedChecks}
            </S.TextDescription>
          </>
        )}
        {!!totalFailedChecks && (
          <>
            <S.CircleCheck cx={failedChecksX} cy={0} r={4} $passed={false} />
            <S.TextDescription dominantBaseline="middle" x={failedChecksX + 6} y={0}>
              {totalFailedChecks}
            </S.TextDescription>
          </>
        )}
      </Group>
      <Group left={outputsX}>
        {hasOutputs && (
          <>
            <S.RectOutput x={0} y={-6} rx={4} />
            <S.TextOutput dominantBaseline="middle" x={2} y={0}>
              O
            </S.TextOutput>
          </>
        )}
      </Group>
    </>
  );
};

export default Header;
