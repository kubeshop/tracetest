import {ArrowRightOutlined} from '@ant-design/icons';
import pagesIcon from 'assets/pages.svg';
import groupIcon from 'assets/group.svg';
import mediaIcon from 'assets/media.svg';
import {COMMUNITY_SLACK_URL, DOCUMENTATION_URL} from 'constants/Common.constants';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import * as S from './Content.styled';

const Content = () => {
  const {navigate} = useDashboard();

  return (
    <S.Container>
      <S.Body>
        <S.Card>
          <S.Icon src={groupIcon} />
          <div>
            <S.Title>Tests</S.Title>
            <S.Text>
              Haven&apos;t created a test yet? Go to the &apos;Tests&apos; page to kickstart your testing adventure.
            </S.Text>
            <S.Link
              onClick={e => {
                e.preventDefault();
                navigate(`/tests`);
              }}
            >
              {' '}
              Go to tests <ArrowRightOutlined />
            </S.Link>
          </div>
        </S.Card>

        <S.Card>
          <S.Icon src={pagesIcon} />
          <div>
            <S.Title>Documentation</S.Title>
            <S.Text>The ultimate technical resources and step-by-step guides that allows you to quickly start.</S.Text>
            <S.Link target="_blank" href={DOCUMENTATION_URL}>
              {' '}
              View documentation <ArrowRightOutlined />
            </S.Link>
          </div>
        </S.Card>

        <S.Card>
          <S.Icon src={mediaIcon} />
          <div>
            <S.Title>Community</S.Title>
            <S.Text>Check the latest updates and support from the community.</S.Text>
            <S.Link target="_blank" href={COMMUNITY_SLACK_URL}>
              {' '}
              Join our community <ArrowRightOutlined />
            </S.Link>
          </div>
        </S.Card>
      </S.Body>
    </S.Container>
  );
};

export default Content;
