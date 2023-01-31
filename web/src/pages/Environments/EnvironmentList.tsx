import Empty from 'components/Empty';
import Pagination from 'components/Pagination';
import usePagination from 'hooks/usePagination';
import Loading from 'pages/Home/Loading';
import {useGetEnvironmentsQuery} from 'redux/apis/TraceTest.api';
import Environment from 'models/Environment.model';
import * as S from './Environment.styled';
import {EnvironmentCard} from './EnvironmentCard';

interface IProps {
  onDelete(id: string): void;
  onEdit(values: Environment): void;
  query: string;
}

const EnvironmentList = ({onDelete, onEdit, query}: IProps) => {
  const pagination = usePagination<Environment, {query: string}>(useGetEnvironmentsQuery, {query});

  return (
    <Pagination
      emptyComponent={
        <Empty message="You have not created any environments yet. Use the Create button to create your first environment" />
      }
      loadingComponent={<Loading />}
      {...pagination}
    >
      <S.ListContainer>
        {pagination.list?.map(environment => (
          <EnvironmentCard environment={environment} key={environment.name} onDelete={onDelete} onEdit={onEdit} />
        ))}
      </S.ListContainer>
    </Pagination>
  );
};

export default EnvironmentList;
