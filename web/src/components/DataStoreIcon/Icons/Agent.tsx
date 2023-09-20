import {IIconProps} from '../DataStoreIcon';

const Agent = ({color, width = '20', height = '20'}: IIconProps) => (
  <svg width={width} height={height} viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M2.14286 0C0.968071 0 0 0.968071 0 2.14286V17.8571C0 19.0319 0.968071 20 2.14286 20H17.8571C19.0319 20 20 19.0319 20 17.8571V2.14286C20 0.968071 19.0319 0 17.8571 0H2.14286ZM3.54492 5C3.73418 4.99332 3.91834 5.06207 4.05692 5.19113L9.05692 9.83398C9.12883 9.90084 9.18618 9.98181 9.22539 10.0718C9.26461 10.1618 9.28485 10.259 9.28485 10.3571C9.28485 10.4553 9.26461 10.5525 9.22539 10.6425C9.18618 10.7325 9.12883 10.8134 9.05692 10.8803L4.05692 15.5232C3.9181 15.6518 3.73392 15.7201 3.54481 15.713C3.35569 15.7059 3.17711 15.6241 3.04827 15.4855C2.91965 15.3467 2.85137 15.1625 2.85843 14.9734C2.86549 14.7843 2.9473 14.6057 3.08592 14.4769L7.5223 10.3571L3.08592 6.23744C2.9473 6.1086 2.86549 5.93001 2.85843 5.7409C2.85137 5.55179 2.91965 5.36761 3.04827 5.22879C3.11197 5.15991 3.18862 5.10426 3.27383 5.065C3.35905 5.02575 3.45117 5.00365 3.54492 5ZM10.0098 12.8571L16.4383 12.9464C16.5321 12.9477 16.6248 12.9674 16.7109 13.0045C16.7971 13.0416 16.8751 13.0953 16.9406 13.1625C17.006 13.2298 17.0575 13.3092 17.0922 13.3964C17.127 13.4835 17.1441 13.5767 17.1429 13.6705C17.1416 13.7643 17.1218 13.8569 17.0848 13.9431C17.0477 14.0293 16.994 14.1073 16.9268 14.1727C16.8595 14.2381 16.7801 14.2897 16.6929 14.3244C16.6058 14.3591 16.5126 14.3763 16.4188 14.375L9.99023 14.2857C9.89643 14.2844 9.8038 14.2647 9.71763 14.2276C9.63145 14.1905 9.55342 14.1369 9.488 14.0696C9.42257 14.0024 9.37104 13.9229 9.33633 13.8358C9.30162 13.7486 9.28442 13.6555 9.28571 13.5617C9.28699 13.4679 9.30673 13.3752 9.34381 13.2891C9.38088 13.2029 9.43457 13.1249 9.50181 13.0594C9.56904 12.994 9.64851 12.9425 9.73566 12.9078C9.82282 12.8731 9.91596 12.8559 10.0098 12.8571Z"
      fill={color}
    />
  </svg>
);

export default Agent;