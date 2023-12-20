import {Button, Col, Space, Typography} from 'antd';

import icon from 'assets/data-stores.svg';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import * as S from './Home.styled';

interface IProps {
  onSkip(): void;
}

const ConfigCTA = ({onSkip}: IProps) => {
  const {navigate} = useDashboard();

  return (
    <S.ConfigContainer align="middle">
      <Col span={12} offset={6}>
        <S.ConfigContent>
          <S.ConfigIcon alt="Tracing Backend" src={icon} />
          <Typography.Title>Configure your Tracing Backend</Typography.Title>
          <Typography.Text>
            Tracetest utilizes the trace collected by your existing OpenTelemetry compatible Tracing Backend to apply
            assertions against. Do you want to configure this now?
          </Typography.Text>
          <S.ConfigFooter>
            <Space>
              <Button onClick={() => navigate('/settings')} type="primary">
                Setup
              </Button>
              <Button data-cy="dataStores-skip-cta" ghost onClick={onSkip} type="primary">
                Skip
              </Button>
            </Space>
          </S.ConfigFooter>
        </S.ConfigContent>
      </Col>
    </S.ConfigContainer>
  );
};

export default ConfigCTA;
