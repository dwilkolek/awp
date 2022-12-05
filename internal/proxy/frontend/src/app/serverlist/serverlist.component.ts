import { Component, OnInit, ViewChild, ViewEncapsulation } from "@angular/core";
import { ServersService } from "src/app/services/servers/servers.service";
import { FormBuilder, FormGroup } from "@angular/forms";
import { Observable } from "rxjs";
import { MatListOption, MatSelectionList } from "@angular/material/list";
import { SocketlogService } from "../services/logs-service/socketlog/socket.service";

@Component({
  selector: "app-serverlist",
  templateUrl: "./serverlist.component.html",
  styleUrls: ["./serverlist.component.sass"],
  encapsulation: ViewEncapsulation.None,
})
export class ServerlistComponent implements OnInit {
  multSelect = false;
  hosts: string[] = [];
  filteredHosts: Observable<string[]>;

  public readonly serverForm: FormGroup;

  @ViewChild("hostsList") hostsList: MatSelectionList;
  selectedOptions: string[] = ["Area 3"];
  colors: Record<string, string>;

  constructor(
    private serversService: ServersService,
    private socketService: SocketlogService,
    private formBuilder: FormBuilder
  ) {
    this.serverForm = this.formBuilder.group({
      serverName: "",
    });
  }

  ngOnInit(): void {
    this.serversService
      .getServers()
      .subscribe((servers) => (this.hosts = servers.hosts));

    this.socketService.colors.subscribe((data) => (this.colors = data));
  }

  onSelectedOptionChange(event: MatListOption[]) {
    if (!!event[0].value) {
      const userSelection = event[0].value;
      this.socketService.subscribeToLog(userSelection);
    }
  }

  get selectedHosts(): string[] {
    return Object.keys(this.colors);
  }

  remove(host: string) {
    this.socketService.unsubscribeFromLog(host);
    this.hostsList.deselectAll();
  }
}
