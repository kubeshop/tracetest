import ResourceCardActions from 'components/ResourceCard/ResourceCardActions';
import {useTransaction} from 'providers/Transaction/Transaction.provider';
import * as S from './TransactionHeader.styled';

interface IProps {
  onBack(): void;
}

const TransactionHeader = ({onBack}: IProps) => {
  const {transaction, onDelete} = useTransaction();

  return (
    <S.Container>
      <S.Section>
        <S.BackIcon data-cy="transaction-header-back-button" onClick={onBack} />
        <div>
          <S.Title data-cy="transaction-details-name">
            {transaction?.name} ({transaction.version})
          </S.Title>
          <S.Text>{transaction?.description}</S.Text>
        </div>
      </S.Section>
      <S.Section>
        <ResourceCardActions id={transaction?.id || ''} onDelete={() => onDelete(transaction.id, transaction.name)} />
      </S.Section>
    </S.Container>
  );
};

export default TransactionHeader;
