import {useCallback} from 'react';
import {CaretUpFilled} from '@ant-design/icons';
import {Collapse, Space, Tooltip, Typography} from 'antd';
import {useAppDispatch} from 'redux/hooks';
import {selectSpan} from 'redux/slices/Trace.slice';
import LinterResult from 'models/LinterResult.model';
import Trace from 'models/Trace.model';
import Span from 'models/Span.model';
import CollapseIcon from './CollapseIcon';
import LintScore from '../LintScore/LintScore';
import * as S from './AnalyzerResult.styled';

interface IProps {
  plugins: LinterResult['plugins'];
  trace: Trace;
}

function getSpanName(spans: Span[], spanId: string) {
  const span = spans.find(s => s.id === spanId);
  return span?.name ?? '';
}

const Plugins = ({plugins, trace}: IProps) => {
  const dispatch = useAppDispatch();

  const onSpanResultClick = useCallback(
    (spanId: string) => {
      dispatch(selectSpan({spanId}));
    },
    [dispatch]
  );

  return (
    <Collapse expandIcon={({isActive = false}) => <CollapseIcon isCollapsed={isActive} />}>
      {plugins.map(plugin => (
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
  );
};

export default Plugins;
