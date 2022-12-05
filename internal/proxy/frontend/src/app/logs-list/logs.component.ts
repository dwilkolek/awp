import { SocketlogService } from "../services/logs-service/socketlog/socket.service";
import { Component, OnInit } from "@angular/core";
import {
  MatSnackBar,
  MatSnackBarHorizontalPosition,
  MatSnackBarVerticalPosition,
} from "@angular/material/snack-bar";

import { ServerLog } from "../services/logs-service/types";

@Component({
  selector: "app-logs",
  templateUrl: "./logs.component.html",
  styleUrls: ["./logs.component.sass"],
})
export class LogsComponent implements OnInit {
  serverLogs: ServerLog[] = [];
  dateNow: Date;
  snackbarDuration = 1.5;
  errorMessage = "";
  horizontalPosition: MatSnackBarHorizontalPosition = "start";
  verticalPosition: MatSnackBarVerticalPosition = "bottom";
  selectedLog: ServerLog | null;
  prettyRes: String = "";
  prettyReq: String = "";

  constructor(
    private socketService: SocketlogService,
    private _snackbar: MatSnackBar
  ) {
    this.dateNow = new Date();
  }

  openSnackBar(message: string) {
    this._snackbar.open(message, "", {
      horizontalPosition: this.horizontalPosition,
      verticalPosition: this.verticalPosition,
      duration: this.snackbarDuration * 1000,
    });
  }

  ngOnInit(): void {
    this.socketService.observableLogs.subscribe(
      (data) => (this.serverLogs = data.reverse())
    );
    this.socketService.noArchiveLogs.subscribe((data) => {
      if (data.length !== 0) {
        this.openSnackBar(data);
      }
    });
  }

  handleLogSelection(log: ServerLog): void {
    this.selectedLog = log;
    this.prettyReq = this.mapLogReq(log);
    this.prettyRes = this.mapLogRes(log);
  }
  mapLogReq(log: ServerLog): string {
    return (
      "Query params:" +
      log.query +
      "\r\n\r\n" +
      "Request Headers\r\n" +
      Object.entries(log.requestHeaders)
        .map((header) => {
          return `${header[0]}: ${
            header[1].length === 1 ? header[1][0] : header[1]
          }`;
        })
        .join("\r\n")
    );
  }
  mapLogRes(log: ServerLog): string {
    const data = log.response.split("\r\n\r\n");
    return data[0] + "\r\n\r\n" + JSON.stringify(JSON.parse(data[1]), null, 2);
  }
}
