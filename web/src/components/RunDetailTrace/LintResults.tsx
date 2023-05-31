import {CaretUpFilled} from '@ant-design/icons';
import {Link} from 'react-router-dom';
import {Col, Collapse, Row, Space, Tooltip, Typography} from 'antd';
import {useCallback} from 'react';
import LinterResult from 'models/LinterResult.model';
import Span from 'models/Span.model';
import Trace from 'models/Trace.model';
import {useAppDispatch} from 'redux/hooks';
import {selectSpan} from 'redux/slices/Trace.slice';
import * as S from './LintResults.styled';
import LintScore from '../LintScore/LintScore';
import CollapseIcon from './CollapseIcon';

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
      <S.Title level={2}>Linter Results</S.Title>
      <S.Description>
        The Tracetest Linter its a plugin based framework used to analyze Open Telemetry traces to help teams improve
        their instrumentation data, find potential problems and provide tips to fix the problems. If you want to disable
        the linter for all tests, go to the <Link to="/settings">settings page</Link>.
      </S.Description>

      <Row gutter={[16, 16]}>
        <Col span={8} key="avg_result">
          <Tooltip title="Tracetest core system supports linter evaluation as part of the testing capabilities.">
            <S.ScoreContainer>
              <S.Subtitle level={3}>Trace Lint Result</S.Subtitle>{' '}
              <LintScore score={linterResult.score} passed={linterResult.passed} />
            </S.ScoreContainer>
          </Tooltip>
        </Col>
        {linterResult?.plugins?.map(plugin => (
          <Col span={8} key={plugin.name}>
            <Tooltip title={plugin.description}>
              <S.ScoreContainer key={plugin.name}>
                <S.Subtitle level={3}>{plugin.name}</S.Subtitle>
                <LintScore score={plugin.score} passed={plugin.passed} />
              </S.ScoreContainer>
            </Tooltip>
          </Col>
        ))}
      </Row>

      <Collapse expandIcon={({isActive = false}) => <CollapseIcon isCollapsed={isActive} />}>
        {linterResult?.plugins?.map(plugin => (
          <S.PluginPanel
            header={
              <Space>
                <LintScore width="35px" height="35px" score={plugin.score} passed={plugin.passed} />
                <Typography.Text strong>{plugin.name}</Typography.Text>
                <Typography.Text type="secondary">{plugin.description}</Typography.Text>
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
                          <div>
                            {result.errors.map(error => (
                              <div>
                                <Typography.Text>{error}</Typography.Text>
                              </div>
                            ))}
                          </div>
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
