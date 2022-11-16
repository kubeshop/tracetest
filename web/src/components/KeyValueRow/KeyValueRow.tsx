import * as S from './KeyValueRow.styled';

interface IProps {
  keyName: string;
  value: string;
}

const KeyValueRow = ({keyName, value}: IProps) => (
  <S.Container>
    <S.Column>
      <S.Text $opacity={0.6}>Key</S.Text>
      <S.Text $fontWeight="600">{`${keyName} `}</S.Text>
    </S.Column>
    <S.Column>
      <S.Text $opacity={0.6}>Value</S.Text>
      <S.Text $fontWeight="600">{`${value}`}</S.Text>
    </S.Column>
  </S.Container>
);

export default KeyValueRow;
