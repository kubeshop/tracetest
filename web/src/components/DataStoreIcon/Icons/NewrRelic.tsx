import {IIconProps} from '../DataStoreIcon';

const NewRelic = ({color, width = '20', height = '20'}: IIconProps) => {
  return (
    <svg width={width} height={height} viewBox="0 0 18 22" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M13.7454 7.88511V14.1152L8.51001 17.2308V21.1248L17.0201 16.0628V5.9375L13.7454 7.88511Z" fill={color} />
      <path
        d="M8.51013 4.77022L13.7455 7.88461L17.0203 5.937L8.51013 0.875L0 5.937L3.27351 7.88461L8.51013 4.77022Z"
        fill={color}
      />
      <path d="M5.23663 12.9484V19.1784L8.51013 21.1248V11.0008L0 5.9375V9.83272L5.23663 12.9484Z" fill={color} />
    </svg>
  );
};

export default NewRelic;
