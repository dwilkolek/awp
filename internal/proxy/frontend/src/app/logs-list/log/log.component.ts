import { SocketlogService } from "./../../services/logs-service/socketlog/socket.service";
import { Component, Input, ViewEncapsulation } from "@angular/core";
import { ServerLog } from "src/app/services/logs-service/types";

@Component({
  selector: "app-log",
  templateUrl: "./log.component.html",
  styleUrls: ["./log.component.sass"],
  encapsulation: ViewEncapsulation.None,
})
export class LogComponent {
  @Input()
  log: ServerLog;
  color: Record<string, string>;

  constructor(private socketService: SocketlogService) {
    socketService.colors.subscribe((data) => (this.color = data));
  }
}
