import {CopyOutlined} from '@ant-design/icons';
import {useTheme} from 'styled-components';
import {THeader} from 'types/Test.types';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import useCopy from 'hooks/useCopy';
import Highlighted from '../Highlighted';
import * as S from './HeaderRow.styled';

interface IProps {
  header: THeader;
}

const HeaderRow = ({header: {key = '', value = ''}}: IProps) => {
  const copy = useCopy();
  const handleOnClick = () => {
    TestRunAnalyticsService.onTriggerResponseHeaderCopy();
    return copy(value);
  };

  const theme = useTheme();

  return (
    <S.HeaderContainer>
      <S.Header>
        <S.HeaderKey>{key}</S.HeaderKey>
        <S.HeaderValue>
          <Highlighted text={value} highlight="" />
        </S.HeaderValue>
      </S.Header>
      <CopyOutlined style={{color: theme.color.textLight}} onClick={handleOnClick} />
    </S.HeaderContainer>
  );
};

export default HeaderRow;
