import {format, formatDistanceToNowStrict, isValid, parseISO} from 'date-fns';
import dropWhile from 'lodash/dropWhile';
import round from 'lodash/round';

export const ONE_MILLISECOND = 1000 * 1;
const ONE_SECOND = 1000 * ONE_MILLISECOND;
const ONE_MINUTE = 60 * ONE_SECOND;
const ONE_HOUR = 60 * ONE_MINUTE;
const ONE_DAY = 24 * ONE_HOUR;

const UNIT_STEPS: {unit: string; microseconds: number; ofPrevious: number}[] = [
  {unit: 'd', microseconds: ONE_DAY, ofPrevious: 24},
  {unit: 'h', microseconds: ONE_HOUR, ofPrevious: 60},
  {unit: 'm', microseconds: ONE_MINUTE, ofPrevious: 60},
  {unit: 's', microseconds: ONE_SECOND, ofPrevious: 1000},
  {unit: 'ms', microseconds: ONE_MILLISECOND, ofPrevious: 1000},
  {unit: 'μs', microseconds: 1, ofPrevious: 1000},
];

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

  isDefaultDate(date: string) {
    return date === '0001-01-01T00:00:00Z';
  },

  /**
   * Format duration for display.
   *
   * @param {number} duration - microseconds
   * @return {string} formatted duration
   */
  formatDuration(duration: number): string {
    // Drop all units that are too large except the last one
    const [primaryUnit, secondaryUnit] = dropWhile(
      UNIT_STEPS,
      ({microseconds}, index) => index < UNIT_STEPS.length - 1 && microseconds > duration
    );

    if (primaryUnit.ofPrevious === 1000) {
      // If the unit is decimal based, display as a decimal
      return `${round(duration / primaryUnit.microseconds, 2)}${primaryUnit.unit}`;
    }

    const primaryValue = Math.floor(duration / primaryUnit.microseconds);
    const primaryUnitString = `${primaryValue}${primaryUnit.unit}`;
    const secondaryValue = Math.round((duration / secondaryUnit.microseconds) % primaryUnit.ofPrevious);
    const secondaryUnitString = `${secondaryValue}${secondaryUnit.unit}`;
    return secondaryValue === 0 ? primaryUnitString : `${primaryUnitString} ${secondaryUnitString}`;
  },
};

export default Date;
