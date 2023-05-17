import {Badge, Col, Collapse, Row, Space, Typography} from 'antd';
import LinterResult from 'models/LinterResult.model';
import * as S from './LintResults.styled';

interface IProps {
  linterResult: LinterResult;
}

const LintResults = ({linterResult}: IProps) => {
  return (
    <S.Container>
      <S.Title level={2}>Lint results</S.Title>

      <S.ScoreContainer>
        <S.Subtitle level={3}>Trace Analysis Result</S.Subtitle>
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
              <S.Subtitle level={3}>{plugin.name}</S.Subtitle>
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

      <Collapse>
        {linterResult?.plugins?.map(plugin => (
          <Collapse.Panel
            header={
              <Space direction="vertical">
                <Badge status={plugin.passed ? 'success' : 'error'} text={plugin.name} />
                <Typography.Text style={{paddingLeft: 14}}>{plugin.description}</Typography.Text>
              </Space>
            }
            key={plugin.name}
          >
            {plugin.rules.map(rule => (
              <S.RuleContainer key={rule.name}>
                <S.Column>
                  <Typography.Text strong>{rule.name}</Typography.Text>
                  <Typography.Text>{rule.description}</Typography.Text>
                  <Typography.Text>{rule.tips.join(' - ')}</Typography.Text>
                </S.Column>

                <S.Column>
                  {rule?.results?.map((result, resultIndex) => (
                    // eslint-disable-next-line react/no-array-index-key
                    <div key={`${result.spanId}-${resultIndex}`}>
                      {result.passed ? (
                        <Typography.Text type="success">SpanId: {result.spanId}</Typography.Text>
                      ) : (
                        <>
                          <Typography.Text>Severity: {result.severity} - </Typography.Text>
                          <Typography.Text>Error: {result.errors.join(' - ')} - </Typography.Text>
                          <Typography.Text type={result.severity === 'error' ? 'danger' : 'warning'}>
                            SpanId: {result.spanId}
                          </Typography.Text>
                        </>
                      )}
                    </div>
                  ))}
                </S.Column>
              </S.RuleContainer>
            ))}
          </Collapse.Panel>
        ))}
      </Collapse>
    </S.Container>
  );
};

export default LintResults;
