import * as React from 'react';
import Date, {ONE_MILLISECOND} from 'utils/Date';
import * as S from './Ticks.styled';

function getLabels(numTicks: number, startTime: number, endTime: number) {
  const viewingDuration = endTime - startTime;
  const labels = [];

  for (let i = 0; i < numTicks; i += 1) {
    const durationAtTick = startTime + (i / (numTicks - 1)) * viewingDuration;
    labels.push(Date.formatDuration(durationAtTick * ONE_MILLISECOND));
  }

  return labels;
}

interface IProps {
  endTime?: number;
  numTicks: number;
  startTime?: number;
}

const Ticks = ({endTime = 0, numTicks, startTime = 0}: IProps) => {
  const labels = getLabels(numTicks, startTime, endTime);

  return (
    <S.Ticks>
      {labels.map((label, index) => {
        const portion = index / (numTicks - 1);
        return (
          <S.Tick key={portion} style={{left: `${portion * 100}%`}}>
            <S.TickLabel $isEndAnchor={portion >= 1}>{label}</S.TickLabel>
          </S.Tick>
        );
      })}
    </S.Ticks>
  );
};

export default Ticks;
