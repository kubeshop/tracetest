import { escapeRegExp } from 'lodash';
import * as S from './Highlighted.styled';

interface IProps {
  text: string;
  highlight: string;
}

const Highlighted = ({text, highlight}: IProps) => {
  if (!highlight.trim()) return <span>{text}</span>;

  const regex = new RegExp(`(${escapeRegExp(highlight)})`, 'gi');
  const partList = text.split(regex);

  return (
    <span>
      {partList.filter(String).map((part, index) =>
        regex.test(part) ? (
          // eslint-disable-next-line react/no-array-index-key
          <S.Mark key={`${part}-${index}`}>{part}</S.Mark>
        ) : (
          // eslint-disable-next-line react/no-array-index-key
          <span key={`${part}-${index}`}>{part}</span>
        )
      )}
    </span>
  );
};

export default Highlighted;
