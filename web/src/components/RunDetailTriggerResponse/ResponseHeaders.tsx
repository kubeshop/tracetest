import HeaderRow from 'components/HeaderRow';
import {isRunStateFinished} from 'models/TestRun.model';
import {TSpanFlatAttribute} from 'types/Span.types';
import {THeader} from 'types/Test.types';
import {TTestRunState} from 'types/TestRun.types';
import SkeletonResponse from './SkeletonResponse';
import * as S from './RunDetailTriggerResponse.styled';

interface IProps {
  headers?: THeader[];
  state: TTestRunState;
  onCreateTestOutput(attribute: TSpanFlatAttribute): void;
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
}

const ResponseHeaders = ({headers, state, onCreateTestOutput, onCreateTestSpec}: IProps) =>
  isRunStateFinished(state) || !!headers ? (
    <S.HeadersList>
      {headers &&
        headers.map(header => (
          <HeaderRow
            header={header}
            key={header.key}
            onCreateTestOutput={onCreateTestOutput}
            onCreateTestSpec={onCreateTestSpec}
          />
        ))}
    </S.HeadersList>
  ) : (
    <SkeletonResponse />
  );

export default ResponseHeaders;
