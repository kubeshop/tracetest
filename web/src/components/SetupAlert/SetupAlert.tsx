import {Button} from 'antd';
import {useDataStoreConfig} from 'providers/DataStoreConfig/DataStoreConfig.provider';
import {Link} from 'react-router-dom';
import {ConfigMode} from 'types/Config.types';
import {useAppSelector} from '../../redux/hooks';
import UserSelectors from '../../selectors/User.selectors';
import * as S from './SetupAlert.styled';

const SetupAlert = () => {
  const {dataStoreConfig, skipConfigSetup} = useDataStoreConfig();
  const initConfigSetup = useAppSelector(state => UserSelectors.selectUserPreference(state, 'initConfigSetup'));
  const shouldDisplay = Boolean(initConfigSetup) && dataStoreConfig.mode === ConfigMode.NO_TRACING_MODE;

  return shouldDisplay ? (
    <S.Container
      message={
        <S.Message>
          <S.TextBold>No trace data store configured.</S.TextBold>
          <S.Text>Let us know the details of your existing tracing solution so we can gather the trace.</S.Text>
          <Link to="/settings">
            <S.WarningButton>Setup</S.WarningButton>
          </Link>
          <Button color="primary" onClick={skipConfigSetup}>
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
