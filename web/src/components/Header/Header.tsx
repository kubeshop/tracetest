import {Space} from 'antd';
import {Link} from 'react-router-dom';
import Logo from 'assets/Logo.svg';
import * as S from './Header.styled';
import HeaderMenu from './HeaderMenu';
import EnvironmentSelector from '../EnvironmentSelector';
import NoTracingPopover from '../NoTracingPopover/NoTracingPopover';

interface IProps {
  hasEnvironments?: boolean;
  hasLogo?: boolean;
  isNoTracingMode: boolean;
}

const Header = ({hasEnvironments = false, hasLogo = false, isNoTracingMode}: IProps) => (
  <S.Header>
    <div>
      {hasLogo && (
        <Link to="/">
          <S.Logo alt="tracetest logo" data-cy="logo" src={Logo} />
        </Link>
      )}
    </div>

    <Space>
      {isNoTracingMode && <NoTracingPopover />}
      {hasEnvironments && <EnvironmentSelector />}

      <HeaderMenu />
    </Space>
  </S.Header>
);

export default Header;
