import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";

export interface AwsServer {
  version:number,
  hosts:string[]
}


@Injectable({providedIn:'root'})
export class ServersService {

  private serversUrl = "/api/config"

  constructor(private http: HttpClient,) {}

  getServers():Observable<AwsServer>{
    return this.http.get<AwsServer>(this.serversUrl)
  }




}
