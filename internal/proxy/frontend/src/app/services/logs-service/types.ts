import { HttpHeaderResponse, HttpHeaders } from "@angular/common/http";

export interface ServerLog {
  color?: string;
  level: string;
  timestamp: number;
  service: string;
  method: string;
  path: string;
  query: string;
  request: string;
  response: string;
  status: number;
  requestHeaders: { [index: string]: string[] };
  responseHeaders: { [index: string]: string[] };
}
