import {useEffect, useState} from 'react';
import TraceAnalyzerAnalytics from 'services/Analytics/TraceAnalyzer.service';
import {TAnalyzerError} from 'types/TestRun.types';
import AnalyzerErrorsPopover from './AnalyzerErrorsPopover';
import * as S from './AnalyzerErrors.styled';

interface IProps {
  errors: TAnalyzerError[];
  isSelected: boolean;
}

const AnalyzerErrors = ({errors, isSelected}: IProps) => {
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    if (isSelected) setIsOpen(true);
  }, [isSelected]);

  return (
    <>
      <S.ErrorIcon
        onClick={() => {
          TraceAnalyzerAnalytics.onSpanErrorsClick();
          setIsOpen(prev => !prev);
        }}
        $isClickable
      />
      {isOpen && <AnalyzerErrorsPopover errors={errors} onClose={() => setIsOpen(false)} />}
    </>
  );
};

export default AnalyzerErrors;
