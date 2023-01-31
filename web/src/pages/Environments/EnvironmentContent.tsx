import {useCallback, useState} from 'react';

import SearchInput from 'components/SearchInput';
import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import {useDeleteEnvironmentMutation} from 'redux/apis/TraceTest.api';
import EnvironmentsAnalytics from 'services/Analytics/EnvironmentsAnalytics.service';
import Environment from 'models/Environment.model';
import * as S from './Environment.styled';
import EnvironmentList from './EnvironmentList';
import {EnvironmentModal} from './EnvironmentModal';

const {onCreateEnvironmentClick} = EnvironmentsAnalytics;

const EnvironmentContent = () => {
  const [deleteEnvironment] = useDeleteEnvironmentMutation();
  const [query, setQuery] = useState<string>('');
  const [environment, setEnvironment] = useState<Environment | undefined>(undefined);
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const onSearch = useCallback((value: string) => setQuery(value), [setQuery]);
  const {onOpen} = useConfirmationModal();

  const handleOnClickCreate = () => {
    onCreateEnvironmentClick();
    setIsModalOpen(true);
  };

  const handleOnModalClose = () => {
    setEnvironment(undefined);
    setIsModalOpen(false);
  };

  const handleOnEdit = (values: Environment) => {
    setEnvironment(values);
    setIsModalOpen(true);
  };

  const handleOnDelete = (id: string) => {
    onOpen({
      title: `Are you sure you want to delete the environment?`,
      onConfirm: () => deleteEnvironment({environmentId: id}),
    });
  };

  return (
    <S.Wrapper>
      <S.MainHeaderContainer>
        <S.TitleText>All Environments</S.TitleText>
      </S.MainHeaderContainer>

      <S.PageHeader>
        <SearchInput onSearch={onSearch} placeholder="Search environment" />
        <S.ActionContainer>
          <S.CreateEnvironmentButton onClick={handleOnClickCreate} type="primary">
            Create Environment
          </S.CreateEnvironmentButton>
        </S.ActionContainer>
      </S.PageHeader>

      <EnvironmentList onDelete={handleOnDelete} onEdit={handleOnEdit} query={query} />
      <EnvironmentModal isOpen={isModalOpen} onClose={handleOnModalClose} environment={environment} />
    </S.Wrapper>
  );
};

export default EnvironmentContent;
