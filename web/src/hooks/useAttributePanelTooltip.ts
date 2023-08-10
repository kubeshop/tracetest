import {useCallback} from 'react';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import UserSelectors from 'selectors/User.selectors';
import {setUserPreference} from '../redux/slices/User.slice';

const tooltip = 'A certain span contains an attribute and this attribute has a specific value. You can check it here.';

const useAttributePanelTooltip = () => {
  const dispatch = useAppDispatch();
  const showAttributeTooltip = useAppSelector(
    state => UserSelectors.selectUserPreference(state, 'showAttributeTooltip') as boolean
  );

  const onClose = useCallback(() => {
    dispatch(setUserPreference({key: 'showAttributeTooltip', value: false}));
  }, [dispatch]);

  return {isVisible: showAttributeTooltip, tooltip, onClose};
};

export default useAttributePanelTooltip;
