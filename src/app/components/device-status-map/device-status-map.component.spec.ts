import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DeviceStatusMap } from './device-status-map.component';

describe('DeviceStatusMap', () => {
  let component: DeviceStatusMap;
  let fixture: ComponentFixture<DeviceStatusMap>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DeviceStatusMap],
    }).compileComponents();

    fixture = TestBed.createComponent(DeviceStatusMap);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
