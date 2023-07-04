import {Group} from '@visx/group';
import * as S from '../Timeline.styled';

interface IProps {
  isLoading: boolean;
  onSelectAsCurrent(): void;
  positionTop: number;
}

const SelectAsCurrent = ({isLoading, onSelectAsCurrent, positionTop}: IProps) => (
  <Group className="matched" left={100} top={positionTop} onClick={() => !isLoading && onSelectAsCurrent()}>
    <S.RectSelectAsCurrent x={0} y={-6} rx={4} />
    <S.TextOutput dominantBaseline="middle" x={isLoading ? 2 : 5} y={0}>
      {isLoading ? 'Updating selected span' : 'Select as current span'}
    </S.TextOutput>
  </Group>
);

export default SelectAsCurrent;
