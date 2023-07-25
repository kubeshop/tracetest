declare module 'parse-curl' {
  export default function parseCurl(command: string): {
    url: string;
    method: string;
    header: {[props: string]: string};
    body: string;
  };
}
