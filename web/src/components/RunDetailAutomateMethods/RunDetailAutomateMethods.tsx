import {Tabs} from 'antd';
import {useSearchParams} from 'react-router-dom';
import {ResourceType} from 'types/Resource.type';
import * as S from './RunDetailAutomateMethods.styled';

interface IProps {
  resourceType: ResourceType;
  methods?: IMethodProps[];
}

interface IMethodProps {
  id: string;
  label: string;
  children: React.ReactNode;
}

const RunDetailAutomateMethods = ({resourceType, methods = []}: IProps) => {
  const [query, updateQuery] = useSearchParams();

  return (
    <S.Container>
      <S.TitleContainer>
        <S.Title>Running Techniques</S.Title>
        <S.Subtitle>Methods to automate the running of this {resourceType}</S.Subtitle>
      </S.TitleContainer>
      <S.TabsContainer>
        <Tabs
          defaultActiveKey={query.get('tab') || methods[0]?.id}
          size="small"
          onChange={activeKey => {
            updateQuery([['tab', activeKey]]);
          }}
        >
          {methods.map(({id, label, children}) => (
            <Tabs.TabPane key={id} tab={label}>
              {children}
            </Tabs.TabPane>
          ))}
        </Tabs>
      </S.TabsContainer>
    </S.Container>
  );
};

export default RunDetailAutomateMethods;
