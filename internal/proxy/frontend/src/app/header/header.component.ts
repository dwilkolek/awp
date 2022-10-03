import { Component, EventEmitter, Input, Output } from "@angular/core";

@Component({
  selector: "app-header",
  templateUrl: "./header.component.html",
  styleUrls: ["./header.component.sass"],
})
export class HeaderComponent {
  @Output() serverButtonEvent = new EventEmitter();
  @Input() tooltipContent = "";

  constructor() {}

  reloadPage() {
    window.location.reload();
  }

  onServerButtonClick() {
    this.serverButtonEvent.emit();
  }
}
