import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {Dropdown, Menu, Typography} from 'antd';
import {useState} from 'react';

import * as T from 'components/ResourceCard/ResourceCard.styled';
import {TEnvironment} from 'types/Environment.types';
import * as E from './Environment.styled';

interface IProps {
  environment: TEnvironment;
  onDelete(id: string): void;
  onEdit(values: TEnvironment): void;
}

export const EnvironmentCard = ({environment: {description, id, name, values}, onDelete, onEdit}: IProps) => {
  const [isCollapsed, setIsCollapsed] = useState(false);

  return (
    <E.EnvironmentCard $isCollapsed={isCollapsed}>
      <E.InfoContainer onClick={() => setIsCollapsed(preCollapsed => !preCollapsed)}>
        {isCollapsed ? <DownOutlined /> : <RightOutlined />}
        <E.TextContainer>
          <E.NameText>{name}</E.NameText>
        </E.TextContainer>
        <E.TextContainer />
        <E.TextContainer>
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
                  label: <span>Edit</span>,
                  onClick: e => {
                    e.domEvent.stopPropagation();
                    onEdit({id, name, description, values});
                  },
                },
                {
                  key: 'delete',
                  label: <span>Delete</span>,
                  onClick: e => {
                    e.domEvent.stopPropagation();
                    onDelete(id);
                  },
                },
              ]}
            />
          }
          placement="bottomLeft"
          trigger={['click']}
        >
          <span onClick={e => e.stopPropagation()}>
            <T.ActionButton />
          </span>
        </Dropdown>
      </E.InfoContainer>

      {isCollapsed && Boolean(values.length) && (
        <E.ResultListContainer>
          <E.VariablesMainContainer>
            <E.HeaderContainer>
              <E.HeaderText>Key</E.HeaderText>
              <E.HeaderTextRight>Value</E.HeaderTextRight>
            </E.HeaderContainer>
            <E.VariablesContainer>
              {values.map(value => (
                <div key={value.key} style={{display: 'flex', justifyContent: 'space-between', width: '100%'}}>
                  <E.VariablesText>{value.key}</E.VariablesText>
                  <E.VariablesText>{value.value}</E.VariablesText>
                </div>
              ))}
            </E.VariablesContainer>
          </E.VariablesMainContainer>
        </E.ResultListContainer>
      )}

      {isCollapsed && !values.length && (
        <T.EmptyStateContainer>
          <T.EmptyStateIcon />
          <Typography.Text disabled>No Variables</Typography.Text>
        </T.EmptyStateContainer>
      )}
    </E.EnvironmentCard>
  );
};
