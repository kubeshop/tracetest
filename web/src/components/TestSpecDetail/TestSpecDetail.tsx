import {TAssertionResultEntry} from 'models/AssertionResults.model';
import * as S from './TestSpecDetail.styled';
import Content from './Content';

interface IProps {
  isOpen: boolean;
  onClose(): void;
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry, name: string): void;
  onRevert(originalSelector: string): void;
  selectedSpan?: string;
  testSpec?: TAssertionResultEntry;
}

const TestSpecDetail = ({isOpen, onClose, onDelete, onEdit, onRevert, selectedSpan, testSpec}: IProps) => (
  <S.DrawerContainer
    closable={false}
    getContainer={false}
    mask={false}
    onClose={onClose}
    placement="right"
    visible={isOpen}
    width="100%"
    height="100%"
  >
    {testSpec && (
      <Content
        onClose={onClose}
        onDelete={onDelete}
        onEdit={onEdit}
        onRevert={onRevert}
        selectedSpan={selectedSpan}
        testSpec={testSpec}
      />
    )}
  </S.DrawerContainer>
);

export default TestSpecDetail;
