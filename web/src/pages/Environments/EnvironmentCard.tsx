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
        <T.TextContainer>
          <E.NameText>{name}</E.NameText>
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
        <T.ResultListContainer>
          <E.VariablesMainContainer>
            <div style={{display: 'flex', justifyContent: 'space-between', paddingBottom: 8}}>
              <Typography style={{flexBasis: '50%', paddingLeft: 8, fontWeight: 'bold'}}>Key</Typography>
              <Typography style={{flexBasis: '50%', fontWeight: 'bold'}}>Value</Typography>
            </div>
            <E.VariablesContainer>
              {resultList.map(secret => (
                <div
                  key={secret.key + secret.value}
                  style={{display: 'flex', justifyContent: 'space-between', width: '100%'}}
                >
                  <Typography style={{flexBasis: '50%'}}>{secret.key}</Typography>
                  <Typography style={{flexBasis: '50%'}}>{secret.value}</Typography>
                </div>
              ))}
            </E.VariablesContainer>
          </E.VariablesMainContainer>
          {resultList.length === 5 && (
            <E.EnvironmentDetails>
              <E.EnvironmentDetailsLink
                data-cy="environment-details-link"
                onClick={() => {
                  // openDialog(environmentId)
                }}
              >
                Explore all environments details
              </E.EnvironmentDetailsLink>
            </E.EnvironmentDetails>
          )}
        </T.ResultListContainer>
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
