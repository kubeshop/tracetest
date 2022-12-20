import {Button} from 'antd';
import {useDataStoreConfig} from 'providers/DataStoreConfig/DataStoreConfig.provider';
import {Link} from 'react-router-dom';
import * as S from './SetupAlert.styled';

const SetupAlert = () => {
  const {shouldDisplayConfigSetupFromTest, skipConfigSetupFromTest} = useDataStoreConfig();

  return shouldDisplayConfigSetupFromTest ? (
    <S.Container
      message={
        <S.Message>
          <S.TextBold>No trace data store configured.</S.TextBold>
          <S.Text>Let us know the details of your existing tracing solution so we can gather the trace.</S.Text>
          <Link to="/settings">
            <S.WarningButton>Setup</S.WarningButton>
          </Link>
          <Button color="primary" onClick={skipConfigSetupFromTest}>
            Later
          </Button>
        </S.Message>
      }
      type="warning"
      showIcon
      closable
    />
  ) : null;
};

export default SetupAlert;
