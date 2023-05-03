import {IIconProps} from '../DataStoreIcon';

const Honeycomb = ({color, width = '20', height = '20'}: IIconProps) => {
  return (
    <svg width={width} height={height} viewBox="0 0 22 20" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path
        d="M10.8606 13.6255L12.6828 16.8117L10.8606 20H7.23773L5.41748 16.8117L7.23773 13.6255H10.8606Z"
        fill={color}
      />
      <path
        d="M10.8606 5.93311L12.6828 9.1194L10.8606 12.3076H7.23773L5.41748 9.1194L7.23773 5.93311H10.8606Z"
        fill={color}
      />
      <path
        d="M4.44754 10.3296L5.93204 12.9649L4.44754 15.6041H1.48448L0 12.9649L1.48448 10.3296H4.44754Z"
        fill={color}
      />
      <path d="M18.7466 0L21.1029 4.17254L18.7466 8.351H14.0104L11.6541 4.17254L14.0104 0H18.7466Z" fill={color} />
    </svg>
  );
};

export default Honeycomb;
