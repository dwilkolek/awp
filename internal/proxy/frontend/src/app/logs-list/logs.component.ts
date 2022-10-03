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
}
