import { WebSocketSubject } from "rxjs/webSocket";
import { ServerlogService } from "src/app/services/logs-service/serverlog/serverlog.service";
import { Injectable } from "@angular/core";
import { BehaviorSubject } from "rxjs";
import { ServerLog } from "../types";

@Injectable({
  providedIn: "root",
})
export class SocketlogService {
  serverLogUrl = "ws://localhost:2137/ws";
  private _observableLogs = new BehaviorSubject<ServerLog[]>([]);
  observableLogs = this._observableLogs.asObservable();
  private _noArchiveLogs = new BehaviorSubject<string>("");
  noArchiveLogs = this._noArchiveLogs.asObservable();
  private _colors = new BehaviorSubject<Record<string, string>>({});
  colors = this._colors.asObservable();

  private socket = new WebSocketSubject({
    url: this.serverLogUrl,
    serializer: (message: string) => message,
  });

  constructor(private serverLog: ServerlogService) {}
  subscribeToLog(host: string) {
    const currentLogs = this._observableLogs.getValue();
    this.serverLog.getLogs(host).subscribe((data) => {
      if (data.length === 0) {
        this._noArchiveLogs.next(`No archive logs available for ${host}`);
      }
      this._observableLogs.next([...currentLogs, ...data]);
      this.socket.next(`sub:${host}`);
      this.socket.asObservable().subscribe((socketLog) => {
        this._observableLogs.next([
          ...this._observableLogs.getValue(),
          socketLog as unknown as ServerLog,
        ]);
      });
    });
  }

  assignColor(host: string) {
    this._colors.next({
      ...this._colors.getValue(),
      [host]: this.generateLightColorHex(),
    });
  }

  unsubscribeFromLog(host: string) {
    this.socket.next(`del:${host}`);
    const filtered = this._observableLogs
      .getValue()
      .filter((server) => server.service !== host);
    this._observableLogs.next([...filtered]);
  }

  private generateLightColorHex() {
    let color = "#";
    for (let i = 0; i < 3; i++)
      color += (
        "0" +
        Math.floor(((1 + Math.random()) * Math.pow(16, 2)) / 2).toString(16)
      ).slice(-2);
    return color;
  }
}
