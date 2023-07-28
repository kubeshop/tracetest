import Empty from 'components/Empty';
import Pagination from 'components/Pagination';
import usePagination from 'hooks/usePagination';
import Loading from 'pages/Home/Loading';
import {useGetEnvironmentsQuery} from 'redux/apis/Tracetest';
import Environment from 'models/Environment.model';
import {ENVIRONMENTS_DOCUMENTATION_URL} from 'constants/Common.constants';
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
        <Empty
          title="You have not created any environments yet"
          message={
            <>
              Use the Create button to create your first environment. Learn more about test or transactions{' '}
              <a href={ENVIRONMENTS_DOCUMENTATION_URL}>here.</a>
            </>
          }
        />
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
