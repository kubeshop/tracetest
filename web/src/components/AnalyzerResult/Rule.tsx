import {useCallback} from 'react';
import {CaretUpFilled} from '@ant-design/icons';
import {Space, Tooltip, Typography} from 'antd';
import {LinterResultPluginRule} from 'models/LinterResult.model';
import Trace from 'models/Trace.model';
import Span from 'models/Span.model';
import {useAppDispatch} from 'redux/hooks';
import {selectSpan} from 'redux/slices/Trace.slice';
import * as S from './AnalyzerResult.styled';
import RuleIcon from './RuleIcon';

interface IProps {
  rule: LinterResultPluginRule;
  trace: Trace;
}

function getSpanName(spans: Span[], spanId: string) {
  const span = spans.find(s => s.id === spanId);
  return span?.name ?? '';
}

const Rule = ({rule: {tips, passed, description, name, errorDescription, results = [], level, weight = 0}, trace}: IProps) => {
  const dispatch = useAppDispatch();

  const onSpanResultClick = useCallback(
    (spanId: string) => {
      dispatch(selectSpan({spanId}));
    },
    [dispatch]
  );

  return (
    <S.RuleContainer>
      <S.Column>
        <S.RuleHeader>
          <Space>
            <RuleIcon passed={passed} level={level} />
            <Tooltip title={tips.join(' - ')}>
              <Typography.Text strong>{name}</Typography.Text>
              {' - '}{weight}%
            </Tooltip>
          </Space>
        </S.RuleHeader>
        <Typography.Text type="secondary" style={{paddingLeft: 20}}>
          {description}
        </Typography.Text>
      </S.Column>

      <S.RuleBody>
        {results?.map((result, resultIndex) => (
          // eslint-disable-next-line react/no-array-index-key
          <div key={`${result.spanId}-${resultIndex}`}>
            <S.SpanButton
              icon={<CaretUpFilled />}
              onClick={() => onSpanResultClick(result.spanId)}
              type="link"
              $error={!result.passed}
            >
              {getSpanName(trace.spans, result.spanId)}
            </S.SpanButton>

            {!result.passed && result.errors.length > 1 && (
              <>
                <div>
                  <Typography.Text>{errorDescription}</Typography.Text>
                </div>
                <S.List>
                  {result.errors.map(error => (
                    <li key={error.value}>
                      <Tooltip title={error.description}>
                        <Typography.Text>{error.value}</Typography.Text>
                      </Tooltip>
                    </li>
                  ))}
                </S.List>
              </>
            )}

            {!result.passed && result.errors.length === 1 && (
              <div>
                <Typography.Text>{result.errors[0].description}</Typography.Text>
              </div>
            )}
          </div>
        ))}
      </S.RuleBody>
    </S.RuleContainer>
  );
};

export default Rule;
