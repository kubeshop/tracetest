import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {Dropdown, Menu, Popover, Typography} from 'antd';
import {useState} from 'react';

import * as T from 'components/ResourceCard/ResourceCard.styled';
import {useFileViewerModal} from 'components/FileViewerModal/FileViewerModal.provider';
import {ResourceType} from 'types/Resource.type';
import Environment from 'models/Environment.model';
import {useEnvironment} from 'providers/Environment/Environment.provider';
import * as E from './Environment.styled';

interface IProps {
  environment: Environment;
  onDelete(id: string): void;
  onEdit(values: Environment): void;
}

export const EnvironmentCard = ({
  environment: {description, id, name, values},
  environment,
  onDelete,
  onEdit,
}: IProps) => {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const {onDefinition} = useFileViewerModal();
  const {selectedEnvironment} = useEnvironment();
  const isSelected = selectedEnvironment?.id === id;

  return (
    <E.EnvironmentCard $isCollapsed={isCollapsed}>
      <E.InfoContainer onClick={() => setIsCollapsed(preCollapsed => !preCollapsed)}>
        {isCollapsed ? <DownOutlined /> : <RightOutlined />}
        <E.NameContainer>
          <E.NameText>{name}</E.NameText>
          {isSelected && (
            <Popover content="Currently selected environment" placement="right">
              <E.InfoIcon />
            </Popover>
          )}
        </E.NameContainer>
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
                  label: <span data-cy="environment-actions-edit">Edit</span>,
                  onClick: e => {
                    e.domEvent.stopPropagation();
                    onEdit(environment);
                  },
                },
                {
                  key: 'definition',
                  label: <span data-cy="environment-actions-definition">Environment Definition</span>,
                  onClick: e => {
                    e.domEvent.stopPropagation();
                    onDefinition(ResourceType.Environment, id);
                  },
                },
                {
                  key: 'delete',
                  label: <span data-cy="environment-actions-delete">Delete</span>,
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
          <span data-cy="environment-actions" onClick={e => e.stopPropagation()}>
            <T.ActionButton />
          </span>
        </Dropdown>
      </E.InfoContainer>

      {isCollapsed && Boolean(values.length) && (
        <E.ResultListContainer>
          <E.VariablesMainContainer>
            <E.HeaderContainer>
              <E.HeaderText>Key</E.HeaderText>
              <E.HeaderText>Value</E.HeaderText>
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
