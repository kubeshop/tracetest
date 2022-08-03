import {useEffect} from 'react';
import {useSpan} from '../../../../providers/Span/Span.provider';
import useControls from '../DAG/hooks/useControls';

interface IProps {
  isSelected?: boolean;
  id: string;
  index: number;
}

export function useScrollSpanNodeGroupIntoView({index, isSelected, id}: IProps): void {
  const {onSelectSpan, affectedSpans} = useSpan();
  const {indexOfFocused} = useControls({onSelectSpan});

  const hasAffectedSpans = affectedSpans.length > 0;
  useEffect(() => {
    if (indexOfFocused === index) {
      if (hasAffectedSpans) {
        document.getElementById(id)?.scrollIntoView();
      }
    }
  }, [id, isSelected, hasAffectedSpans, indexOfFocused, index]);
}
