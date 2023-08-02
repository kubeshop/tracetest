import {useState} from 'react';

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
        onClick={e => {
          e.stopPropagation();
          setIsOpen(true);
        }}
        $isClickable
      />
      {isOpen && <AnalyzerErrorsPopover errors={errors} onClose={() => setIsOpen(false)} />}
    </>
  );
};

export default AnalyzerErrors;
