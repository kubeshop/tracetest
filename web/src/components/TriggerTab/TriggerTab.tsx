import {useTheme} from 'styled-components';
import * as S from './TriggerTab.styled';

interface IProps {
  hasContent?: boolean;
  label: string;
  totalItems?: number;
}

const TriggerTab = ({hasContent, label, totalItems}: IProps) => {
  const {
    color: {primary},
  } = useTheme();

  return (
    <div>
      {`${label} `}
      {!!totalItems && <S.Text>({totalItems})</S.Text>}
      {hasContent && <S.Badge color={primary} />}
    </div>
  );
};

export default TriggerTab;
