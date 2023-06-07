import {ClockCircleOutlined, SettingOutlined, ToolOutlined} from '@ant-design/icons';
import {useMemo} from 'react';
import {Handle, NodeProps, Position} from 'react-flow-renderer';

import AssertionResultChecks from 'components/AssertionResultChecks/AssertionResultChecks';
import CurrentSpanSelector from 'components/CurrentSpanSelector';
import TestOutputMark from 'components/TestOutputMark';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import {SemanticGroupNamesToText} from 'constants/SemanticGroupNames.constants';
import {SpanKindToText} from 'constants/Span.constants';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {selectAnalyzerResults} from 'redux/slices/Trace.slice';
import {selectOutputsBySpanId} from 'redux/testOutputs/selectors';
import TestSpecsSelectors from 'selectors/TestSpecs.selectors';
import TraceSelectors from 'selectors/Trace.selectors';
import SpanService from 'services/Span.service';
import {INodeDataSpan} from 'types/DAG.types';
import AnalyzerResults from './AnalyzerResults';
import * as SA from './AnalyzerResults.styled';
import * as S from './SpanNode.styled';

interface IProps extends NodeProps<INodeDataSpan> {}

const SpanNode = ({data, id, selected}: IProps) => {
  const dispatch = useAppDispatch();
  const assertions = useAppSelector(state => TestSpecsSelectors.selectAssertionResultsBySpan(state, data?.id || ''));
  const outputs = useAppSelector(state => selectOutputsBySpanId(state, data?.id || ''));
  const {failed, passed} = useMemo(() => SpanService.getAssertionResultSummary(assertions), [assertions]);
  const {isOpen: isTestSpecFormOpen} = useTestSpecForm();
  const {isOpen: isTestOutputFormOpen} = useTestOutput();
  const {matchedSpans} = useSpan();
  const {runLinterResultsBySpan} = useTestRun();
  const lintErrors = useMemo(
    () => SpanService.filterLintErrorsBySpan(runLinterResultsBySpan, data.id),
    [runLinterResultsBySpan, data.id]
  );
  const showSelectAsCurrent =
    !data.isMatched && !!matchedSpans.length && (isTestSpecFormOpen || isTestOutputFormOpen) && selected;
  const className = `${data.isMatched ? 'matched' : ''} ${showSelectAsCurrent ? 'selectedAsCurrent' : ''}`;
  const selectedAnalyzerResults = useAppSelector(TraceSelectors.selectSelectedAnalyzerResults);

  const handleSelectAnalyzerResults = (spanId: string = '') => {
    dispatch(selectAnalyzerResults({spanId}));
  };

  return (
    <>
      <S.Container
        className={className}
        data-cy={`trace-node-${data.type}`}
        $matched={data.isMatched}
        $selected={selected}
      >
        <Handle id={id} position={Position.Top} style={{top: 0, visibility: 'hidden'}} type="target" />

        <S.TopLine $type={data.type} />

        {selectedAnalyzerResults === data.id && (
          <AnalyzerResults lintErrors={lintErrors} onClose={() => handleSelectAnalyzerResults()} />
        )}

        <S.Header>
          <S.BadgeContainer>
            <S.BadgeType count={SemanticGroupNamesToText[data.type]} $hasMargin $type={data.type} />
          </S.BadgeContainer>
          <S.HeaderText>{data.name}</S.HeaderText>
          {!!lintErrors.length && <SA.ErrorIcon $isAbsolute onClick={() => handleSelectAnalyzerResults(data.id)} />}
        </S.Header>

        <S.Body>
          <S.Item>
            <SettingOutlined />
            <S.ItemText>
              {data.service} {SpanKindToText[data.kind]}
            </S.ItemText>
          </S.Item>
          {Boolean(data.system) && (
            <S.Item>
              <ToolOutlined />
              <S.ItemText>{data.system}</S.ItemText>
            </S.Item>
          )}
          <S.Item>
            <ClockCircleOutlined />
            <S.ItemText>{data.duration}</S.ItemText>
          </S.Item>
        </S.Body>

        <S.Footer>
          {!!outputs.length && <TestOutputMark outputs={outputs} />}
          <AssertionResultChecks failed={failed} passed={passed} styleType="node" />
        </S.Footer>

        <Handle id={id} position={Position.Bottom} style={{bottom: 0, visibility: 'hidden'}} type="source" />
      </S.Container>
      {showSelectAsCurrent && <CurrentSpanSelector spanId={data.id} />}
    </>
  );
};

export default SpanNode;
