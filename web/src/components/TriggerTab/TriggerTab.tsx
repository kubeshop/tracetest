import * as S from './TriggerTab.styled';

interface IProps {
  hasContent?: boolean;
  label: string;
  totalItems?: number;
}

const TriggerTab = ({hasContent, label, totalItems}: IProps) => {
  return (
    <div>
      {`${label} `}
      {!!totalItems && <S.Text>({totalItems})</S.Text>}
      {hasContent && <S.Badge status="success" />}
    </div>
  );
};

export default TriggerTab;
