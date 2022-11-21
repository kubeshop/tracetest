import {DownOutlined} from '@ant-design/icons';
import {Dropdown, Menu, Space} from 'antd';
import {useEnvironment} from 'providers/Environment/Environment.provider';

const EnvironmentSelector = () => {
  const {environmentList, selectedEnvironment, setSelectedEnvironment, isLoading} = useEnvironment();

  const menu = (
    <Menu
      items={[
        {
          key: 'no-env',
          label: 'No environment',
          onClick: () => setSelectedEnvironment(),
        },
      ].concat(
        environmentList.map(environment => ({
          key: environment.id,
          label: environment.name,
          onClick: () => setSelectedEnvironment(environment),
        }))
      )}
    />
  );

  return !isLoading ? (
    <Dropdown overlay={menu} overlayClassName="environment-selector-items">
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
