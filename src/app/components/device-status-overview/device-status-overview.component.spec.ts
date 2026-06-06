import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DeviceStatusOverview } from './device-status-overview.component';

describe('DeviceStatusOverview', () => {
  let component: DeviceStatusOverview;
  let fixture: ComponentFixture<DeviceStatusOverview>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DeviceStatusOverview],
    }).compileComponents();

    fixture = TestBed.createComponent(DeviceStatusOverview);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
