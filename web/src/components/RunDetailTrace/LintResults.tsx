import {CaretUpFilled} from '@ant-design/icons';
import {Col, Collapse, Row, Space, Tooltip, Typography} from 'antd';
import {useCallback} from 'react';
import {TooltipQuestion} from 'components/TooltipQuestion/TooltipQuestion';
import LinterResult from 'models/LinterResult.model';
import Span from 'models/Span.model';
import Trace from 'models/Trace.model';
import {useAppDispatch} from 'redux/hooks';
import {selectSpan} from 'redux/slices/Trace.slice';
import * as S from './LintResults.styled';

interface IProps {
  linterResult: LinterResult;
  trace: Trace;
}

function getSpanName(spans: Span[], traceId: string) {
  const span = spans.find(s => s.id === traceId);
  return span?.name ?? '';
}

const LintResults = ({linterResult, trace}: IProps) => {
  const dispatch = useAppDispatch();

  const onSpanResultClick = useCallback(
    (spanId: string) => {
      dispatch(selectSpan({spanId}));
    },
    [dispatch]
  );

  return (
    <S.Container>
      <S.Title level={2}>Lint results</S.Title>

      <S.ScoreContainer>
        <S.Subtitle level={3}>
          Trace Analysis Result
          <TooltipQuestion title="Tracetest core system supports linter evaluation as part of the testing capabilities." />
        </S.Subtitle>{' '}
        <Space>
          <S.Score level={1}>{linterResult.score} %</S.Score>
          <S.ScoreProgress
            format={() => ''}
            percent={linterResult.score}
            status={linterResult.passed ? 'success' : 'exception'}
            type="circle"
          />
        </Space>
      </S.ScoreContainer>

      <Row gutter={[16, 16]}>
        {linterResult?.plugins?.map(plugin => (
          <Col span={12} key={plugin.name}>
            <S.ScoreContainer key={plugin.name}>
              <S.Subtitle level={3}>
                {plugin.name}
                <TooltipQuestion title={plugin.description} />
              </S.Subtitle>
              <Space>
                <S.Score level={1}>{plugin.score} %</S.Score>
                <S.ScoreProgress
                  format={() => ''}
                  percent={plugin.score}
                  status={plugin.passed ? 'success' : 'exception'}
                  type="circle"
                />
              </Space>
            </S.ScoreContainer>
          </Col>
        ))}
      </Row>

      <Collapse expandIcon={() => null}>
        {linterResult?.plugins?.map(plugin => (
          <S.PluginPanel
            header={
              <Space>
                {plugin.passed ? <S.PassedIcon /> : <S.FailedIcon />}
                <Typography.Text strong>{plugin.name}</Typography.Text>
                <TooltipQuestion title={plugin.description} />
              </Space>
            }
            key={plugin.name}
          >
            {plugin.rules.map(rule => (
              <S.RuleContainer key={rule.name}>
                <S.Column>
                  <S.RuleHeader>
                    <Space>
                      {rule.passed ? <S.PassedIcon $small /> : <S.FailedIcon $small />}
                      <Tooltip title={rule.tips.join(' - ')}>
                        <Typography.Text strong>{rule.name}</Typography.Text>
                      </Tooltip>
                    </Space>
                  </S.RuleHeader>
                  <Typography.Text type="secondary" style={{paddingLeft: 20}}>
                    {rule.description}
                  </Typography.Text>
                </S.Column>

                <S.RuleBody>
                  {rule?.results?.map((result, resultIndex) => (
                    // eslint-disable-next-line react/no-array-index-key
                    <div key={`${result.spanId}-${resultIndex}`}>
                      {result.passed ? (
                        <S.SpanButton
                          icon={<CaretUpFilled />}
                          onClick={() => onSpanResultClick(result.spanId)}
                          type="link"
                        >
                          {getSpanName(trace.spans, result.spanId)}
                        </S.SpanButton>
                      ) : (
                        <>
                          <S.SpanButton
                            icon={<CaretUpFilled />}
                            onClick={() => onSpanResultClick(result.spanId)}
                            type="link"
                            $error
                          >
                            {getSpanName(trace.spans, result.spanId)}
                          </S.SpanButton>
                          <Typography.Text> - Severity: {result.severity}</Typography.Text>
                          <Typography.Text> - Errors: {result.errors.join(' - ')}</Typography.Text>
                        </>
                      )}
                    </div>
                  ))}
                </S.RuleBody>
              </S.RuleContainer>
            ))}
          </S.PluginPanel>
        ))}
      </Collapse>
    </S.Container>
  );
};

export default LintResults;
