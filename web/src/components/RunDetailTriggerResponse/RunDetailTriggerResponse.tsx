import {Tabs} from 'antd';
import {TriggerTypes} from '../../constants/Test.constants';
import {TTriggerResult} from '../../types/Test.types';
import HeaderRow from './HeaderRow';
import * as S from './RunDetailTriggerResponse.styled';

interface IProps {
  triggerResult?: TTriggerResult;
  executionTime?: number;
}

const RunDetailTriggerResponse = ({
  executionTime = 0,
  triggerResult: {headers = [], body = '{}', statusCode = 200} = {
    headers: [],
    body: '{}',
    type: TriggerTypes.http,
    statusCode: 200,
  },
}: IProps) => {
  const onCopy = (value: string) => {
    navigator.clipboard.writeText(value);
  };

  return (
    <S.Container>
      <S.TitleContainer>
        <S.Title>Response Data</S.Title>
        <div>
          <S.StatusText>
            Status: <S.StatusSpan $isError={statusCode >= 400}>{statusCode}</S.StatusSpan>
          </S.StatusText>
          <S.StatusText>
            Time: <S.StatusSpan $isError={executionTime > 1000}>{executionTime}ms</S.StatusSpan>
          </S.StatusText>
        </div>
      </S.TitleContainer>
      <S.TabsContainer>
        <Tabs defaultActiveKey="1" data-cy="run-detail-trigger-response" size="small">
          <Tabs.TabPane key="1" tab="Body">
            <S.ValueJson>
              <pre>{JSON.stringify(JSON.parse(body), null, 2)}</pre>
            </S.ValueJson>
          </Tabs.TabPane>
          <Tabs.TabPane key="2" tab="Headers">
            {headers.map(header => (
              <HeaderRow onCopy={onCopy} header={header} key={header.key} />
            ))}
          </Tabs.TabPane>
        </Tabs>
      </S.TabsContainer>
    </S.Container>
  );
};

export default RunDetailTriggerResponse;
