import {COMMUNITY_SLACK_URL, GITHUB_ISSUES_URL} from 'constants/Common.constants';
import * as S from './ContactUs.styled';

const ContactUsModalFooter = () => {
  return (
    <S.ModalFooter>
      <a href={GITHUB_ISSUES_URL} target="_blank">
        <S.FullWidthButton type="primary">Create an Issue</S.FullWidthButton>
      </a>
      <a href={COMMUNITY_SLACK_URL} target="_blank">
        <S.FullWidthButton ghost type="primary">
          Contact Team on Slack
        </S.FullWidthButton>
      </a>
    </S.ModalFooter>
  );
};

export default ContactUsModalFooter;
