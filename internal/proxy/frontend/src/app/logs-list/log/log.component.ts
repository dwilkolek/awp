import { SocketlogService } from "./../../services/logs-service/socketlog/socket.service";
import {
  Component,
  EventEmitter,
  Input,
  Output,
  ViewEncapsulation,
} from "@angular/core";
import { ServerLog } from "src/app/services/logs-service/types";
import { NgClass } from "@angular/common";

@Component({
  selector: "app-log",
  templateUrl: "./log.component.html",
  styleUrls: ["./log.component.sass"],
  encapsulation: ViewEncapsulation.None,
})
export class LogComponent {
  @Input()
  log: ServerLog;

  @Input()
  selected: boolean = false;

  @Output()
  clicked: EventEmitter<void> = new EventEmitter();

  color: Record<string, string>;

  constructor(private socketService: SocketlogService) {
    socketService.colors.subscribe((data) => (this.color = data));
  }

  public select() {
    this.clicked.emit();
  }
}
