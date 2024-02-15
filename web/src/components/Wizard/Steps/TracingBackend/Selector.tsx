import {Typography} from 'antd';
import {COMMUNITY_SLACK_URL} from 'constants/Common.constants';
import {SupportedDataStores} from 'types/DataStore.types';
import * as S from './TracingBackend.styled';
import BackendCard from './BackendCard';

const backends = Object.values(SupportedDataStores);

interface IProps {
  onSelect(backend: SupportedDataStores): void;
  selectedBackend?: SupportedDataStores;
}

const Selector = ({onSelect, selectedBackend}: IProps) => {
  return (
    <S.Container>
      <Typography.Title level={1}>Tell us how Tracetest should ingest traces from your application</Typography.Title>

      <Typography.Paragraph type="secondary">Select your Tracing Backend</Typography.Paragraph>

      <S.BackendSelector>
        {backends.map(backend => (
          <BackendCard
            backend={backend}
            selectedBackend={selectedBackend}
            key={backend}
            onSelect={() => onSelect(backend)}
          />
        ))}
      </S.BackendSelector>

      <Typography.Paragraph type="secondary">
        Don&apos;t see the Tracing Backend you need?{' '}
        <a href={COMMUNITY_SLACK_URL} target="__blank">
          Submit a request.
        </a>
      </Typography.Paragraph>
    </S.Container>
  );
};

export default Selector;
