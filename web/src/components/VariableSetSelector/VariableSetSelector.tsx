import {DownOutlined} from '@ant-design/icons';
import {Dropdown, Menu, Space} from 'antd';
import type {ItemType} from 'antd/lib/menu/hooks/useItems';
import {Operation, useCustomization} from 'providers/Customization';
import {useVariableSet} from 'providers/VariableSet';
import {useMemo, useState} from 'react';
import AddVariableSet from './AddVariableSet';
import VariableSetSelectorEntry from './VariableSetSelectorEntry';

const VariableSetSelector = () => {
  const {getIsAllowed} = useCustomization();
  const {variableSetList, selectedVariableSet, setSelectedVariableSet, isLoading, onOpenModal} = useVariableSet();
  const [hoveredOption, setHoveredOption] = useState<string>();

  const items = useMemo<ItemType[]>(
    () =>
      (
        [
          {
            key: 'no-variable-set',
            label: 'No Variable Set',
            onClick: () => setSelectedVariableSet(),
          },
        ] as ItemType[]
      )
        .concat(
          variableSetList.map(variableSet => ({
            key: variableSet.id,
            label: (
              <VariableSetSelectorEntry
                onEditClick={onOpenModal}
                variableSet={variableSet}
                isHovering={hoveredOption === variableSet.id}
                isAllowed={getIsAllowed(Operation.Edit)}
              />
            ),
            onClick: () => setSelectedVariableSet(variableSet),
            onMouseEnter: () => setHoveredOption(variableSet.id),
          }))
        )
        .concat([
          {
            type: 'divider',
          },
          {
            key: 'add-variable-set',
            label: <AddVariableSet />,
            onClick: () => onOpenModal(),
            disabled: !getIsAllowed(Operation.Edit),
          },
        ]),
    [getIsAllowed, hoveredOption, onOpenModal, setSelectedVariableSet, variableSetList]
  );

  return !isLoading ? (
    <Dropdown overlay={<Menu items={items} />} overlayClassName="variableSet-selector-items">
      <a data-cy="variableSet-selector" onClick={e => e.preventDefault()}>
        <Space>
          {selectedVariableSet?.name || 'No Variable Set'}
          <DownOutlined />
        </Space>
      </a>
    </Dropdown>
  ) : null;
};

export default VariableSetSelector;
