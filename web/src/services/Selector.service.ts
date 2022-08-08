import {tracetestLang} from '../utils/grammar';

const SelectorService = () => ({
  getIsValidSelector(query: string): boolean {
    try {
      tracetestLang.parser.configure({strict: true}).parse(query);
      return true;
    } catch (e) {
      return false;
    }
  },
});

export default SelectorService();
