import {CaretUpFilled} from '@ant-design/icons';
import {Space, Switch, Tooltip, Typography} from 'antd';
import {useCallback, useState} from 'react';
import AnalyzerScore from 'components/AnalyzerScore/AnalyzerScore';
import LinterResult from 'models/LinterResult.model';
import Trace from 'models/Trace.model';
import Span from 'models/Span.model';
import {useAppDispatch} from 'redux/hooks';
import {selectSpan} from 'redux/slices/Trace.slice';
import AnalyzerService from 'services/Analyzer.service';
import * as S from './AnalyzerResult.styled';
import CollapseIcon from './CollapseIcon';

interface IProps {
  plugins: LinterResult['plugins'];
  trace: Trace;
}

function getSpanName(spans: Span[], spanId: string) {
  const span = spans.find(s => s.id === spanId);
  return span?.name ?? '';
}

const Plugins = ({plugins: rawPlugins, trace}: IProps) => {
  const dispatch = useAppDispatch();
  const [onlyErrors, setOnlyErrors] = useState(false);
  const plugins = AnalyzerService.getPlugins(rawPlugins, onlyErrors);

  const onSpanResultClick = useCallback(
    (spanId: string) => {
      dispatch(selectSpan({spanId}));
    },
    [dispatch]
  );

  return (
    <>
      <S.SwitchContainer>
        <Switch checked={onlyErrors} id="only_errors_enabled" onChange={() => setOnlyErrors(prev => !prev)} />
        <label htmlFor="only_errors_enabled">Show only errors</label>
      </S.SwitchContainer>

      <S.StyledCollapse expandIcon={({isActive = false}) => <CollapseIcon isCollapsed={isActive} />}>
        {plugins.map(plugin => (
          <S.PluginPanel
            header={
              <Space>
                <AnalyzerScore width="35px" height="35px" score={plugin.score} />
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
                            {result.groupedErrors.map((groupedError, index) => (
                              // eslint-disable-next-line react/no-array-index-key
                              <div key={index}>
                                <div>
                                  <Typography.Text>{groupedError.error}</Typography.Text>
                                </div>
                                <S.List>
                                  {groupedError.values?.map(value => (
                                    <li key={value}>
                                      <Typography.Text>{value}</Typography.Text>
                                    </li>
                                  ))}
                                </S.List>
                              </div>
                            ))}

                            {!result.groupedErrors.length &&
                              result.errors.map((error, index) => (
                                // eslint-disable-next-line react/no-array-index-key
                                <div key={index}>
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
      </S.StyledCollapse>
    </>
  );
};

export default Plugins;
