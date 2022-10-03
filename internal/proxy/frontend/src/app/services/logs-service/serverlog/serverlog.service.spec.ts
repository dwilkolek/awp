import { TestBed } from "@angular/core/testing";

import { ServerlogService } from "src/app/services/logs-service/serverlog/serverlog.service";

describe("ServerlogService", () => {
  let service: ServerlogService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ServerlogService);
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });
});
