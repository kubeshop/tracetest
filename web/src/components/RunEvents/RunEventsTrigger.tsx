import {Tabs, Typography} from 'antd';
import {useSearchParams} from 'react-router-dom';
import {COMMUNITY_SLACK_URL, GITHUB_ISSUES_URL} from 'constants/Common.constants';
import RunEvent from './RunEvent';
import {IPropsComponent} from './RunEvents';
import * as S from './RunEvents.styled';
import ResponseVariableSet from '../RunDetailTriggerResponse/ResponseVariableSet';
import ResponseMetadata from '../RunDetailTriggerResponse/ResponseMetadata';

const TabsKeys = {
  Error: 'error',
  VariableSet: 'variable-set',
  Metadata: 'metadata',
};

const RunEventsTrigger = ({events}: IPropsComponent) => {
  const [query, updateQuery] = useSearchParams();

  return (
    <>
      <S.TitleContainer>
        <S.Title>Trigger Error</S.Title>
      </S.TitleContainer>
      <S.TabsContainer>
        <Tabs
          defaultActiveKey={query.get('tab') || TabsKeys.Error}
          data-cy="run-detail-trigger-response"
          size="small"
          onChange={newTab => {
            updateQuery([['tab', newTab]]);
          }}
        >
          <Tabs.TabPane key={TabsKeys.Error} tab="Event Log">
            <S.Container>
              <S.ErrorIcon />
              <Typography.Title level={2}>Test Trigger Failed</Typography.Title>
              <S.Paragraph>
                We prepared the breakdown of diagnostic steps and tips to help identify the source of the issue:
              </S.Paragraph>

              <S.ListContainer>
                {events.map(event => (
                  <RunEvent event={event} key={event.type} />
                ))}
              </S.ListContainer>

              <S.Paragraph type="secondary">
                <S.Link href={GITHUB_ISSUES_URL} target="_blank">
                  Create an issue
                </S.Link>{' '}
                or contact us via{' '}
                <S.Link href={COMMUNITY_SLACK_URL} target="_blank">
                  Slack
                </S.Link>
                . We will check it out and will help you rectify the issue.
              </S.Paragraph>
            </S.Container>
          </Tabs.TabPane>
          <Tabs.TabPane key={TabsKeys.VariableSet} tab="Variable Set">
            <ResponseVariableSet />
          </Tabs.TabPane>
          <Tabs.TabPane key={TabsKeys.Metadata} tab="Metadata">
            <ResponseMetadata />
          </Tabs.TabPane>
        </Tabs>
      </S.TabsContainer>
    </>
  );
};

export default RunEventsTrigger;
