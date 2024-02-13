import {Tooltip, Typography} from 'antd';
import {CaretUpFilled} from '@ant-design/icons';
import {useCallback, useMemo} from 'react';
import {LinterResultPluginRule} from 'models/LinterResult.model';
import {useAppDispatch} from 'redux/hooks';
import {selectSpan} from 'redux/slices/Trace.slice';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import TraceAnalyzerAnalytics from 'services/Analytics/TraceAnalyzer.service';
import * as S from './AnalyzerResult.styled';
import RuleLink from './RuleLink';

interface IProps {
  index: number;
  data: LinterResultPluginRule;
  style: React.CSSProperties;
}

const RuleResult = ({index, data: {results, id, errorDescription}, style}: IProps) => {
  const {spanId, passed, errors} = useMemo(() => results[index], [results, index]);
  const dispatch = useAppDispatch();
  const {
    run: {trace},
  } = useTestRun();

  const onClick = useCallback(() => {
    TraceAnalyzerAnalytics.onSpanNameClick();
    dispatch(selectSpan({spanId}));
  }, [dispatch, spanId]);

  return (
    <div key={`${spanId}-${index}`} style={style}>
      <S.SpanButton icon={<CaretUpFilled />} onClick={onClick} type="link" $error={!passed}>
        {trace.flat[spanId].name ?? ''}
      </S.SpanButton>

      {!passed && errors.length > 1 && (
        <>
          <div>
            <Typography.Text>{errorDescription}</Typography.Text>
          </div>
          <S.List>
            {errors.map(error => (
              <li key={error.value}>
                <Tooltip title={error.description}>
                  <Typography.Text>{error.value}</Typography.Text>
                </Tooltip>
              </li>
            ))}
          </S.List>
        </>
      )}

      {!passed && errors.length === 1 && (
        <div>
          <Typography.Text>{errors[0].description}</Typography.Text>
        </div>
      )}

      {!passed && <RuleLink id={id} />}
    </div>
  );
};

export default RuleResult;
