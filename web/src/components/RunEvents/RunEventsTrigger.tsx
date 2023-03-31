import {Typography} from 'antd';

import {DISCORD_URL, GITHUB_ISSUES_URL} from 'constants/Common.constants';
import RunEvent from './RunEvent';
import {IPropsComponent} from './RunEvents';
import * as S from './RunEvents.styled';

const RunEventsTrigger = ({events}: IPropsComponent) => (
  <S.Container>
    <S.ErrorIcon />
    <Typography.Title level={2}>Test Run Failed</Typography.Title>
    <S.Paragraph>
      We prepared the breakdown of diagnostic steps and tips to help identify the source of the issue:
    </S.Paragraph>

    <S.ListContainer>
      {events.map(event => (
        <RunEvent event={event} key={event.type} />
      ))}
    </S.ListContainer>

    <S.Paragraph>
      <Typography.Text type="secondary">
        <S.Link href={GITHUB_ISSUES_URL} target="_blank">
          Create an issue
        </S.Link>{' '}
        or contact us via{' '}
        <S.Link href={DISCORD_URL} target="_blank">
          Discord
        </S.Link>
        . We will check it out and will help you rectify the issue.
      </Typography.Text>
    </S.Paragraph>
  </S.Container>
);

export default RunEventsTrigger;
