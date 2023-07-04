import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';

interface IProps {
  selected: boolean;
  matched: boolean;
}

const useShowSelectAsCurrent = ({selected, matched}: IProps): boolean => {
  const {matchedSpans} = useSpan();
  const {isOpen: isTestSpecFormOpen} = useTestSpecForm();
  const {isOpen: isTestOutputFormOpen} = useTestOutput();

  const showSelectAsCurrent =
    selected && !matched && !!matchedSpans.length && (isTestSpecFormOpen || isTestOutputFormOpen);

  return showSelectAsCurrent;
};

export default useShowSelectAsCurrent;
