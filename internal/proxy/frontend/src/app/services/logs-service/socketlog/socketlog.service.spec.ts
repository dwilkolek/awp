import { TestBed } from "@angular/core/testing";
import { SocketlogService } from "./socket.service";

describe("SocketlogService", () => {
  let service: SocketlogService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(SocketlogService);
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });
});
