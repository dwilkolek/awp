import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { catchError, map, Observable, retry, throwError } from "rxjs";
import { ServerLog } from "../types";

@Injectable({
  providedIn: "root",
})
export class ServerlogService {
  serverLogUrl = "api/logs";

  constructor(private http: HttpClient) {}

  getLogs(serverName: string): Observable<ServerLog[]> {
    return this.http
      .get<ServerLog[]>(`${this.serverLogUrl}/${serverName}`)
      .pipe(
        map((data) => {
          if (!data) {
            return [];
          }
          return data;
        }),
        catchError((error) => {
          console.log(error);
          throw throwError(() => new Error("Error fetching logs"));
        }),
        retry(2)
      );
  }
}
