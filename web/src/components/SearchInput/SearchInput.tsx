import {debounce} from 'lodash';
import {useMemo} from 'react';
import * as S from './SearchInput.styled';

interface ISearchInputProps {
  height?: string;
  width?: string;
  placeholder: string;
  onSearch(value: string): void;
  delay?: number;
}

const SearchInput: React.FC<ISearchInputProps> = ({
  height = '32px',
  width = '270px',
  placeholder,
  onSearch,
  delay = 500,
}) => {
  const handleSearch = useMemo(
    () =>
      debounce(event => {
        onSearch(event.target.value);
      }, delay),
    [delay, onSearch]
  );

  return (
    <S.SearchInput
      prefix={<S.SearchIcon />}
      placeholder={placeholder}
      width={width}
      height={height}
      onChange={handleSearch}
      allowClear
    />
  );
};

export default SearchInput;
