import {useEffect, useState} from 'react';

const useDraftMode = (isDefaultDraftMode = false) => {
  const [isDraftMode, setIsDraftMode] = useState(isDefaultDraftMode);

  useEffect(() => {
    if (isDraftMode) window.onbeforeunload = () => true;
    else window.onbeforeunload = null;
  }, [isDraftMode]);

  return {isDraftMode, setIsDraftMode};
};

export default useDraftMode;
