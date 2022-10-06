import {CopyOutlined} from '@ant-design/icons';
import {message} from 'antd';
import {useTheme} from 'styled-components';
import {THeader} from 'types/Test.types';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import Highlighted from '../Highlighted';
import * as S from './HeaderRow.styled';

interface IProps {
  header: THeader;
  onCopy(value: string): void;
}

const HeaderRow = ({header: {key = '', value = ''}, onCopy}: IProps) => {
  const handleOnClick = () => {
    message.success('Value copied to the clipboard');
    TestRunAnalyticsService.onTriggerResponseHeaderCopy();
    return onCopy(value);
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
