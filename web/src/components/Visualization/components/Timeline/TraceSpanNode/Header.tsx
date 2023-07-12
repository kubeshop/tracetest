import errorIcon from 'assets/error.svg';
import * as S from '../Timeline.styled';

interface IProps {
  hasAnalyzerErrors: boolean;
}

const Header = ({hasAnalyzerErrors}: IProps) => {
  if (!hasAnalyzerErrors) return null;

  return (
    <>
      <S.Image href={errorIcon} x={0} y={-5} $height={10} $width={10} />
      <S.TextDescription dominantBaseline="middle" x={14} y={0}>
        Analyzer errors
      </S.TextDescription>
    </>
  );
};

export default Header;
