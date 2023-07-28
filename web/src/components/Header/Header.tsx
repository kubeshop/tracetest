import {Space} from 'antd';

import Logo from 'assets/Logo.svg';
import EnvironmentSelector from 'components/EnvironmentSelector';
import Link from 'components/Link';
import NoTracingPopover from 'components/NoTracingPopover';
import * as S from './Header.styled';
import HelpMenu from './HelpMenu';

interface IProps {
  hasLogo?: boolean;
  isNoTracingMode: boolean;
}

const Header = ({hasLogo = false, isNoTracingMode}: IProps) => (
  <S.Header>
    <div>
      {hasLogo && (
        <Link to="/">
          <S.Logo alt="Tracetest logo" data-cy="logo" src={Logo} />
        </Link>
      )}
    </div>

    <Space>
      {isNoTracingMode && <NoTracingPopover />}
      <EnvironmentSelector />
      <HelpMenu />
    </Space>
  </S.Header>
);

export default Header;
