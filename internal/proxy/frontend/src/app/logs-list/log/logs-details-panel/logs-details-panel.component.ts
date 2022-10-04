import { Component, Input } from "@angular/core";
import { ServerLog } from "../../../services/logs-service/types";

@Component({
  selector: "app-logs-details-panel",
  templateUrl: "./logs-details-panel.component.html",
  styleUrls: ["./logs-details-panel.component.sass"],
})
export class LogsDetailsPanelComponent {
  panelOpenState = false;
  @Input() log: ServerLog;

  constructor() {}
}
