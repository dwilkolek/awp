import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LogsDetailsPanelComponent } from './logs-details-panel.component';

describe('LogsDetailsPanelComponent', () => {
  let component: LogsDetailsPanelComponent;
  let fixture: ComponentFixture<LogsDetailsPanelComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ LogsDetailsPanelComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LogsDetailsPanelComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
