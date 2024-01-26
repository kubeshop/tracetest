import {useTheme} from 'styled-components';
import {SupportedIntegrations} from 'constants/Integrations.constants';
import Cypress from './Icons/Cypress';
import K6 from './Icons/K6';
import Playwright from './Icons/Playwright';

const iconMap = {
  [SupportedIntegrations.CYPRESS]: Cypress,
  [SupportedIntegrations.PLAYWRIGHT]: Playwright,
  [SupportedIntegrations.K6]: K6,
} as const;

interface IProps {
  color?: string;
  integrationName: SupportedIntegrations;
  width?: string;
  height?: string;
}

export interface IIconProps {
  color: string;
  width?: string;
  height?: string;
}

const IntegrationIcon = ({color, integrationName, width, height}: IProps) => {
  const {
    color: {text: defaultColor},
  } = useTheme();
  const Component = iconMap[integrationName];

  return <Component color={color || defaultColor} width={width} height={height} />;
};

export default IntegrationIcon;
