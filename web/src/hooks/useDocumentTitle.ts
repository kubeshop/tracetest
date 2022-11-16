import {useEffect} from 'react';
import {DOCUMENT_TITLE} from '../constants/Common.constants';

const useDocumentTitle = (title: string) => {
  useEffect(() => {
    document.title = `${DOCUMENT_TITLE} - ${title}`;

    return () => {
      document.title = DOCUMENT_TITLE;
    };
  }, [title]);
};

export default useDocumentTitle;
