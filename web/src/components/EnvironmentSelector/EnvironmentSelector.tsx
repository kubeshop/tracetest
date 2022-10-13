import {DownOutlined} from '@ant-design/icons';
import {Dropdown, Menu, Space} from 'antd';
import {useEnvironment} from 'providers/Environment/Environment.provider';

const EnvironmentSelector = () => {
  const {environmentList, selectedEnvironment, setSelectedEnvironment} = useEnvironment();

  const menu = (
    <Menu
      items={environmentList.map(environment => ({
        key: environment.id,
        label: environment.name,
        onClick: () => setSelectedEnvironment(environment),
      }))}
    />
  );

  return (
    <Dropdown overlay={menu}>
      <a onClick={e => e.preventDefault()}>
        <Space>
          {selectedEnvironment?.name || 'All Environments'}
          <DownOutlined />
        </Space>
      </a>
    </Dropdown>
  );
};

export default EnvironmentSelector;
