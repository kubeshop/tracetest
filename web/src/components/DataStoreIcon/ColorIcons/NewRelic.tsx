import {IIconProps} from '../DataStoreIcon';

const NewRelic = ({width = '20', height = '20'}: IIconProps) => (
  <svg width={width} height={height} viewBox="0 0 21 24" fill="none" xmlns="http://www.w3.org/2000/svg">
    <g clipPath="url(#clip0_1083_22079)">
      <path
        d="M16.7836 8.30695V15.6896L10.3911 19.3817V23.996L20.7821 17.9975V5.99902L16.7836 8.30695Z"
        fill="#00AC69" />
      <path
        d="M10.3909 4.61683L16.7834 8.30739L20.7819 5.99947L10.3909 0.000976562L0 5.99947L3.99698 8.30739L10.3909 4.61683Z"
        fill="#1CE783" />
      <path d="M6.39396 14.3069V21.6896L10.3909 23.996V11.999L0 5.99902V10.6149L6.39396 14.3069Z" fill="#1D252C" />
    </g>
    <defs>
      <clipPath id="clip0_1083_22079">
        <rect width="20.8302" height="24" fill="white" transform="translate(0 0.000976562)" />
      </clipPath>
    </defs>
  </svg>
);

export default NewRelic;
