import * as S from './ContactUs.styled';
import ContactUsModalFooter from './ContactUsModalFooter';

interface IProps {
  isOpen: boolean;
  onClose(): void;
}

const ContactUsModal = ({isOpen, onClose}: IProps) => {
  return (
    <S.Modal visible={isOpen} onCancel={onClose} footer={<ContactUsModalFooter />}>
      <S.Header>
        <S.PlushieImage width="100px" height="auto" />
        <S.Title>Let us help you</S.Title>
      </S.Header>
      <S.Message>
        Technical glitches can be tricky, even for the best of us. Don&apos;t fret! We&apos;re here to save your day.
        Create an Issue or contact us via Slack and our tech-savvy team are ready to lend a helping hand.
      </S.Message>
    </S.Modal>
  );
};

export default ContactUsModal;
