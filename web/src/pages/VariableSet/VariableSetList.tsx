import Empty from 'components/Empty';
import Pagination from 'components/Pagination';
import usePagination from 'hooks/usePagination';
import Loading from 'pages/Home/Loading';
import TracetestAPI from 'redux/apis/Tracetest';
import VariableSet from 'models/VariableSet.model';
import {VARIABLE_SET_DOCUMENTATION_URL} from 'constants/Common.constants';
import * as S from './VariableSet.styled';
import VariableSetCard from './VariableSetCard';

const {useGetVariableSetsQuery} = TracetestAPI.instance;

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
          title="Haven't Created a VariableSet Yet"
          message={
            <>
              Hit the &apos;Create&apos; button to create your first variable set. Want to learn more about
              VariableSets? Just click{' '}
              <S.Link href={VARIABLE_SET_DOCUMENTATION_URL} target="_blank">
                here.
              </S.Link>
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
