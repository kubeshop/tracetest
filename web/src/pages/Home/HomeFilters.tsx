import {useCallback} from 'react';
import {Select, Typography} from 'antd';
import SearchInput from 'components/SearchInput';
import {SortBy, SortDirection, sortOptions} from 'constants/Test.constants';
import * as S from './Home.styled';

interface IProps {
  onSearch(search: string): void;
  onSortBy(sortBy: SortBy, sortDirection: SortDirection): void;
  isEmpty: boolean;
}

const HomeFilters = ({onSearch, onSortBy, isEmpty}: IProps) => {
  const handleSort = useCallback(
    (newSortBy: string) => {
      const {
        params: {sortBy, sortDirection},
      } = sortOptions.find(({value}) => value === newSortBy)!;

      onSortBy(sortBy, sortDirection);
    },
    [onSortBy]
  );

  return (
    <S.FiltersContainer>
      <SearchInput onSearch={onSearch} placeholder="Search test" />
      <Typography.Text>Sort by:</Typography.Text>
      <Select disabled={isEmpty} defaultValue={sortOptions[0].value} onChange={handleSort} style={{width: 160}}>
        {sortOptions.map(({value, label}) => (
          <Select.Option key={value} value={value}>
            {label}
          </Select.Option>
        ))}
      </Select>
    </S.FiltersContainer>
  );
};

export default HomeFilters;
