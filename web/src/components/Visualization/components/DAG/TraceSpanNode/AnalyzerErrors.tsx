import {useState} from 'react';
import TraceAnalyzerAnalytics from 'services/Analytics/TraceAnalyzer.service';
import {TAnalyzerError} from 'types/TestRun.types';
import AnalyzerErrorsPopover from './AnalyzerErrorsPopover';
import * as S from './AnalyzerErrors.styled';

interface IProps {
  errors: TAnalyzerError[];
}

const AnalyzerErrors = ({errors}: IProps) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <S.ErrorIcon
        onClick={() => {
          TraceAnalyzerAnalytics.onSpanErrorsClick();
          setIsOpen(true);
        }}
        $isClickable
      />
      {isOpen && <AnalyzerErrorsPopover errors={errors} onClose={() => setIsOpen(false)} />}
    </>
  );
};

export default AnalyzerErrors;
