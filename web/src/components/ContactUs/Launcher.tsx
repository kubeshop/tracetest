import {withCustomization} from 'providers/Customization';
import * as S from './ContactUs.styled';
import PulseButton from '../PulseButton';

interface IProps {
  onClick(): void;
}

const Launcher = ({onClick}: IProps) => (
  <S.Container onClick={onClick}>
    <S.PulseButtonContainer>
      <PulseButton />
    </S.PulseButtonContainer>
    <S.PlushieImage />
  </S.Container>
);

export default withCustomization(Launcher, 'contactLauncher');
