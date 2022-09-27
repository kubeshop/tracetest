import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {Dropdown, Menu, Typography} from 'antd';
import {useCallback, useState} from 'react';
import * as T from '../../components/TestCard/TestCard.styled';
import {IEnvironment, useLazyGetEnvironmentSecretListQuery} from '../../redux/apis/TraceTest.api';
import EnvironmentsAnalytics from '../../services/Analytics/EnvironmentsAnalytics.service';
import * as S from './Envs.styled';

interface IProps {
  openDialog: (mode: boolean) => void;
  environment: IEnvironment;
  setEnvironment: (mode: IEnvironment) => void;
}

export const EnvironmentCard = ({
  setEnvironment,
  openDialog,
  environment: {name, id, description},
}: IProps): React.ReactElement => {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [loadResultList, {data: resultList = []}] = useLazyGetEnvironmentSecretListQuery();
  const onCollapse = useCallback(async () => {
    EnvironmentsAnalytics.onEnvironmentClick(id);

    if (resultList.length > 0) {
      setIsCollapsed(true);
      return;
    }
    await loadResultList({environmentId: id, take: 5});
    setIsCollapsed(true);
  }, [loadResultList, resultList.length, id]);

  const toggleColapsed = useCallback(() => {
    if (isCollapsed) {
      setIsCollapsed(false);
      return;
    }
    onCollapse();
  }, [onCollapse, isCollapsed, setIsCollapsed]);
  return (
    <T.TestCard $isCollapsed={isCollapsed}>
      <T.InfoContainer onClick={toggleColapsed}>
        {isCollapsed ? <DownOutlined /> : <RightOutlined data-cy={`collapse-environment-${id}`} />}
        <T.TextContainer>
          <T.NameText>{name}</T.NameText>
        </T.TextContainer>
        <T.TextContainer />
        <T.TextContainer data-cy={`environment-description-${id}`}>
          <T.Text>{description}</T.Text>
        </T.TextContainer>
        <T.TextContainer />
        <T.TextContainer />

        <Dropdown
          overlay={
            <Menu
              items={[
                {
                  key: 'edit',
                  label: <span data-cy="environment-card-edit">Edit</span>,
                  onClick: e => {
                    e.domEvent.stopPropagation();
                    setEnvironment({id, name, description, variables: []});
                    openDialog(true);
                  },
                },
              ]}
            />
          }
          placement="bottomLeft"
          trigger={['click']}
        >
          <span
            data-cy={`environment-actions-button-${id}`}
            className="ant-dropdown-link"
            onClick={e => e.stopPropagation()}
          >
            <T.ActionButton />
          </span>
        </Dropdown>
      </T.InfoContainer>

      {isCollapsed && Boolean(resultList.length) && (
        <T.ResultListContainer>
          <S.VariablesMainContainer>
            <div style={{display: 'flex', justifyContent: 'space-between', paddingBottom: 8}}>
              <Typography style={{flexBasis: '50%', paddingLeft: 8, fontWeight: 'bold'}}>Key</Typography>
              <Typography style={{flexBasis: '50%', fontWeight: 'bold'}}>Value</Typography>
            </div>
            <S.VariablesContainer>
              {resultList.map(secret => (
                <div style={{display: 'flex', justifyContent: 'space-between', width: '100%'}}>
                  <Typography style={{flexBasis: '50%'}}>{secret.key}</Typography>
                  <Typography style={{flexBasis: '50%'}}>{secret.value}</Typography>
                </div>
              ))}
            </S.VariablesContainer>
          </S.VariablesMainContainer>
          {resultList.length === 5 && (
            <T.TestDetails>
              <T.TestDetailsLink
                data-cy="test-details-link"
                onClick={() => {
                  // openDialog(environmentId)
                }}
              >
                Explore all test details
              </T.TestDetailsLink>
            </T.TestDetails>
          )}
        </T.ResultListContainer>
      )}

      {isCollapsed && !resultList.length && (
        <T.EmptyStateContainer>
          <T.EmptyStateIcon />
          <Typography.Text disabled>No Variables</Typography.Text>
        </T.EmptyStateContainer>
      )}
    </T.TestCard>
  );
};
