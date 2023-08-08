import Empty from 'components/Empty';
import Pagination from 'components/Pagination';
import usePagination from 'hooks/usePagination';
import Loading from 'pages/Home/Loading';
import {useGetVariableSetsQuery} from 'redux/apis/Tracetest';
import VariableSet from 'models/VariableSet.model';
import {VARIABLE_SET_DOCUMENTATION_URL} from 'constants/Common.constants';
import * as S from './VariableSet.styled';
import VariableSetCard from './VariableSetCard';

interface IProps {
  onDelete(id: string): void;
  onEdit(values: VariableSet): void;
  query: string;
}

const VariableSetList = ({onDelete, onEdit, query}: IProps) => {
  const pagination = usePagination<VariableSet, {query: string}>(useGetVariableSetsQuery, {query});

  return (
    <Pagination
      emptyComponent={
        <Empty
          title="You have not created any variable sets yet"
          message={
            <>
              Use the Create button to create your first variable set. Learn more about test or test suites{' '}
              <a href={VARIABLE_SET_DOCUMENTATION_URL}>here.</a>
            </>
          }
        />
      }
      loadingComponent={<Loading />}
      {...pagination}
    >
      <S.ListContainer>
        {pagination.list?.map(variableSet => (
          <VariableSetCard variableSet={variableSet} key={variableSet.name} onDelete={onDelete} onEdit={onEdit} />
        ))}
      </S.ListContainer>
    </Pagination>
  );
};

export default VariableSetList;
