import {DownOutlined} from '@ant-design/icons';
import {Dropdown, Menu, Space} from 'antd';
import type {ItemType} from 'antd/lib/menu/hooks/useItems';
import {useEnvironment} from 'providers/Environment/Environment.provider';
import {useMemo, useState} from 'react';
import AddEnvironment from './AddEnvironment';
import EnvironmentSelectorEntry from './EnvironmentSelectorEntry';

const EnvironmentSelector = () => {
  const {environmentList, selectedEnvironment, setSelectedEnvironment, isLoading} = useEnvironment();
  const [hoveredOption, setHoveredOption] = useState<string>();
  const {onOpenModal} = useEnvironment();

  const items = useMemo<ItemType[]>(
    () =>
      (
        [
          {
            key: 'no-env',
            label: 'No environment',
            onClick: () => setSelectedEnvironment(),
          },
        ] as ItemType[]
      )
        .concat(
          environmentList.map(environment => ({
            key: environment.id,
            label: (
              <EnvironmentSelectorEntry
                onEditClick={onOpenModal}
                environment={environment}
                isHovering={hoveredOption === environment.id}
              />
            ),
            onClick: () => setSelectedEnvironment(environment),
            onMouseEnter: () => setHoveredOption(environment.id),
          }))
        )
        .concat([
          {
            type: 'divider',
          },
          {
            key: 'add-env',
            label: <AddEnvironment />,
            onClick: () => onOpenModal(),
          },
        ]),
    [environmentList, hoveredOption, onOpenModal, setSelectedEnvironment]
  );

  return !isLoading ? (
    <Dropdown overlay={<Menu items={items} />} overlayClassName="environment-selector-items">
      <a data-cy="environment-selector" onClick={e => e.preventDefault()}>
        <Space>
          {selectedEnvironment?.name || 'No environment'}
          <DownOutlined />
        </Space>
      </a>
    </Dropdown>
  ) : null;
};

export default EnvironmentSelector;
