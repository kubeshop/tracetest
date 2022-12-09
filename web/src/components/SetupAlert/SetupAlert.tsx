import {useDataStoreConfig} from 'providers/DataStoreConfig/DataStoreConfig.provider';
import {ConfigMode} from 'types/Config.types';
import * as S from './SetupAlert.styled';

const SetupAlert = () => {
  const {dataStoreConfig} = useDataStoreConfig();
  const shouldDisplay = dataStoreConfig.mode === ConfigMode.NO_TRACING_MODE;

  return shouldDisplay ? (
    <S.Container
      message={
        <S.Message>
          <S.TextBold>No trace data store configured.</S.TextBold>
          <S.Text>Let us know the details of your existing tracing solution so we can gather the trace.</S.Text>
          <S.WarningButton href="/settings">Setup</S.WarningButton>
        </S.Message>
      }
      type="warning"
      showIcon
      closable
    />
  ) : null;
};

export default SetupAlert;
