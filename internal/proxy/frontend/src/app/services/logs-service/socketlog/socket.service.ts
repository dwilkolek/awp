import { WebSocketSubject } from "rxjs/webSocket";
import { ServerlogService } from "src/app/services/logs-service/serverlog/serverlog.service";
import { Injectable } from "@angular/core";
import { BehaviorSubject } from "rxjs";
import { ServerLog } from "../types";

@Injectable({
  providedIn: "root",
})
export class SocketlogService {
  subscriptionColorMap: { [index: string]: string } = JSON.parse(
    localStorage.getItem("subscriptionColorMap") || "{}"
  );
  serverLogUrl = "ws://localhost:2137/ws";
  private logs: ServerLog[] = [];
  private _observableLogs = new BehaviorSubject<ServerLog[]>([]);
  observableLogs = this._observableLogs.asObservable();

  private _noArchiveLogs = new BehaviorSubject<string>("");
  noArchiveLogs = this._noArchiveLogs.asObservable();

  private _colors = new BehaviorSubject<Record<string, string>>(
    this.subscriptionColorMap
  );
  colors = this._colors.asObservable();

  private socket = new WebSocketSubject({
    url: this.serverLogUrl,
    serializer: (message: string) => message,
  });

  constructor(private serverLog: ServerlogService) {
    Object.entries(this.subscriptionColorMap).forEach((entry) =>
      this.subscribeToLog(entry[0], entry[1])
    );
  }

  subscribeToLog(host: string, color: string | undefined = undefined) {
    console.log("subsribing to ", host, color);

    this.serverLog.getLogs(host).subscribe((data) => {
      if (data.length === 0) {
        this._noArchiveLogs.next(`No archive logs available for ${host}`);
      }
      this.logs = [...this.logs, ...data];
      this._observableLogs.next(this.sortedLogsCopy());
      this.socket.next(`sub:${host}`);
      this.socket.asObservable().subscribe((socketLog) => {
        this.logs.push(socketLog as unknown as ServerLog);
        this._observableLogs.next(this.sortedLogsCopy());
      });
    });
    if (!color) {
      this.assignColor(host);
    }
  }

  private sortedLogsCopy() {
    return [...this.logs.sort((a, b) => a.timestamp - b.timestamp)];
  }

  assignColor(host: string) {
    const nextState = {
      ...this._colors.getValue(),
      [host]: this.generateLightColorHex(),
    };
    this._colors.next(nextState);
    localStorage.setItem("subscriptionColorMap", JSON.stringify(nextState));
  }

  unsubscribeFromLog(host: string) {
    this.socket.next(`del:${host}`);
    this.logs = this.logs.filter((server) => server.service !== host);

    const nextState = {
      ...this._colors.getValue(),
    };
    delete nextState[host];
    this._colors.next(nextState);
    localStorage.setItem("subscriptionColorMap", JSON.stringify(nextState));
    this._observableLogs.next(this.sortedLogsCopy());
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
