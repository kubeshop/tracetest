import {TAssertionResultEntry} from 'types/Assertion.types';
import * as S from './TestSpecDetail.styled';
import Content from './Content';

interface IProps {
  isOpen: boolean;
  onClose(): void;
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry): void;
  onRevert(originalSelector: string): void;
  onSelectSpan(spanId: string): void;
  selectedSpan?: string;
  testSpec?: TAssertionResultEntry;
}

const TestSpecDetail = ({
  isOpen,
  onClose,
  onDelete,
  onEdit,
  onRevert,
  onSelectSpan,
  selectedSpan,
  testSpec,
}: IProps) => {
  return (
    <S.DrawerContainer
      closable={false}
      getContainer={false}
      mask={false}
      onClose={onClose}
      placement="right"
      visible={isOpen}
      width="100%"
    >
      {testSpec && (
        <Content
          onClose={onClose}
          onDelete={onDelete}
          onEdit={onEdit}
          onRevert={onRevert}
          onSelectSpan={onSelectSpan}
          selectedSpan={selectedSpan}
          testSpec={testSpec}
        />
      )}
    </S.DrawerContainer>
  );
};

export default TestSpecDetail;
