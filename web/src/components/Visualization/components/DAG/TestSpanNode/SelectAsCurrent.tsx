import {LoadingOutlined} from '@ant-design/icons';
import * as S from './TestSpanNode.styled';

interface IProps {
  isLoading: boolean;
  onSelectAsCurrent(): void;
}

const SelectAsCurrent = ({isLoading, onSelectAsCurrent}: IProps) => (
  <S.SelectAsCurrentContainer className="matched" onClick={() => !isLoading && onSelectAsCurrent()}>
    <S.FloatingText>
      {isLoading ? (
        <>
          Updating selected span <LoadingOutlined />
        </>
      ) : (
        <>Select as current span</>
      )}
    </S.FloatingText>
  </S.SelectAsCurrentContainer>
);

export default SelectAsCurrent;
