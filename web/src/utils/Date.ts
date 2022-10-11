import {format, formatDistanceToNowStrict, isValid, parseISO} from 'date-fns';

const Date = {
  format(date: string, dateFormat = "EEEE, yyyy/MM/dd 'at' HH:mm:ss") {
    const isoDate = parseISO(date);
    if (!isValid(isoDate)) {
      return '';
    }
    return format(isoDate, dateFormat);
  },
  getTimeAgo(date: string) {
    const isoDate = parseISO(date);
    if (!isValid(isoDate)) {
      return '';
    }
    return formatDistanceToNowStrict(isoDate, {addSuffix: true});
  },
};

export default Date;
