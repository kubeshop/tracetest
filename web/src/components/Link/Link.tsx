import {AnchorHTMLAttributes, MouseEvent, useCallback, useContext} from 'react';
import {Link as RRLink, LinkProps} from 'react-router-dom';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';

export interface IProps extends LinkProps {}

const Link = ({to, className, children, ...rest}: IProps) => {
  const {baseUrl} = useDashboard();

  let prefixedTo = to;
  if (typeof to === 'string') {
    prefixedTo = `${baseUrl}${to}`;
  } else if (typeof to === 'object') {
    prefixedTo = {...to, pathname: `${baseUrl}${to.pathname}`};
  }

  return (
    <RRLink to={prefixedTo} className={className} {...rest}>
      {children}
    </RRLink>
  );
};

export default Link;
