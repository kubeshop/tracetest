import AttributeActions from 'components/AttributeActions';
import {isRunStateFinished} from 'models/TestRun.model';
import {TSpanFlatAttribute} from 'types/Span.types';
import {TTestRunState} from 'types/TestRun.types';
import SkeletonResponse from './SkeletonResponse';
import CodeBlock from '../CodeBlock';
import * as S from './RunDetailTriggerResponse.styled';

interface IProps {
  body?: string;
  bodyMimeType?: string;
  state: TTestRunState;
  onCreateTestOutput(attribute: TSpanFlatAttribute): void;
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
}

const ResponseBody = ({body = '', bodyMimeType = '', state, onCreateTestOutput, onCreateTestSpec}: IProps) =>
  isRunStateFinished(state) || !!body ? (
    <S.ResponseBodyContainer>
      <S.ResponseBodyContent>
        <CodeBlock value={body} mimeType={bodyMimeType} maxHeight="540px" />
      </S.ResponseBodyContent>
      <S.ResponseBodyActions>
        <AttributeActions
          attribute={{key: 'tracetest.response.body', value: body}}
          onCreateTestOutput={onCreateTestOutput}
          onCreateTestSpec={onCreateTestSpec}
        />
      </S.ResponseBodyActions>
    </S.ResponseBodyContainer>
  ) : (
    <SkeletonResponse />
  );

export default ResponseBody;
