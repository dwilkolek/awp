import { HttpHeaderResponse, HttpHeaders } from "@angular/common/http";

export interface ServerLog {
  color?: string;
  level: string;
  ts: string;
  service: string;
  method: string;
  path: string;
  query: string;
  request: string;
  response: string;
  status: number;
  requestHeaders: HttpHeaders;
  responseHeaders: HttpHeaderResponse;
}
