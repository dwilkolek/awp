import { NgModule } from "@angular/core";
import { BrowserModule } from "@angular/platform-browser";

import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { HeaderComponent } from "./header/header.component";
import { LogsComponent } from "./logs-list/logs.component";
import { HttpClientModule } from "@angular/common/http";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MaterialDesignModule } from "src/material.module";
import { ServerlistComponent } from "./serverlist/serverlist.component";
import { FilterPipe } from "src/app/serverlist/FilterPipe";
import { LogComponent } from "./logs-list/log/log.component";
import { KeyvalueComponent } from "./logs-list/keyvalue/keyvalue.component";
import { LogsDetailsPanelComponent } from "./logs-list/log/logs-details-panel/logs-details-panel.component";
import { ScrollingModule } from "@angular/cdk/scrolling";

@NgModule({
  declarations: [
    AppComponent,
    HeaderComponent,
    LogsComponent,
    ServerlistComponent,
    FilterPipe,
    LogComponent,
    KeyvalueComponent,
    LogsDetailsPanelComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    BrowserAnimationsModule,
    ReactiveFormsModule,
    MaterialDesignModule,
    FormsModule,
    ScrollingModule,
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
