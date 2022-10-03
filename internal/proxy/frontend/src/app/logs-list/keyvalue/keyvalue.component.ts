import { Component, Input } from "@angular/core";

@Component({
  selector: "app-keyvalue",
  templateUrl: "./keyvalue.component.html",
  styleUrls: ["./keyvalue.component.sass"],
})
export class KeyvalueComponent {
  @Input()
  label: string;

  @Input()
  value: string | null | undefined;

  constructor() {}
}
