import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {Dropdown, Menu, Typography} from 'antd';
import {Dispatch, SetStateAction, useCallback, useState} from 'react';
import {useLazyGetEnvironmentSecretListQuery} from 'redux/apis/TraceTest.api';
import * as T from '../../components/TestCard/TestCard.styled';
import EnvironmentsAnalytics from '../../services/Analytics/EnvironmentsAnalytics.service';
import * as E from './Environment.styled';
import {IEnvironment} from './IEnvironment';

interface IProps {
  setIsFormOpen: Dispatch<SetStateAction<boolean>>;
  environment: IEnvironment;
  setEnvironment: Dispatch<SetStateAction<IEnvironment | undefined>>;
}

export const EnvironmentCard = ({
  setIsFormOpen,
  setEnvironment,
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
    <E.EnvironmentCard $isCollapsed={isCollapsed}>
      <E.InfoContainer onClick={toggleColapsed}>
        {isCollapsed ? <DownOutlined /> : <RightOutlined data-cy={`collapse-environment-${id}`} />}
        <E.TextContainer>
          <E.NameText>{name}</E.NameText>
        </E.TextContainer>
        <E.TextContainer />
        <E.TextContainer data-cy={`environment-description-${id}`}>
          <T.Text>{description}</T.Text>
        </E.TextContainer>
        <E.TextContainer />
        <E.TextContainer />

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
                    setIsFormOpen(true);
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
      </E.InfoContainer>

      {isCollapsed && Boolean(resultList.length) && (
        <E.ResultListContainer>
          <E.VariablesMainContainer>
            <E.HeaderContainer>
              <E.HeaderText>Key</E.HeaderText>
              <E.HeaderTextRight>Value</E.HeaderTextRight>
            </E.HeaderContainer>
            <E.VariablesContainer>
              {resultList.map(secret => (
                <div
                  key={secret.key + secret.value}
                  style={{display: 'flex', justifyContent: 'space-between', width: '100%'}}
                >
                  <E.VariablesText>{secret.key}</E.VariablesText>
                  <E.VariablesText>{secret.value}</E.VariablesText>
                </div>
              ))}
            </E.VariablesContainer>
          </E.VariablesMainContainer>
          {resultList.length === 5 && (
            <E.EnvironmentDetails>
              <E.EnvironmentDetailsLink data-cy="environment-details-link">
                Explore all environments details
              </E.EnvironmentDetailsLink>
            </E.EnvironmentDetails>
          )}
        </E.ResultListContainer>
      )}

      {isCollapsed && !resultList.length && (
        <T.EmptyStateContainer>
          <T.EmptyStateIcon />
          <Typography.Text disabled>No Variables</Typography.Text>
        </T.EmptyStateContainer>
      )}
    </E.EnvironmentCard>
  );
};
