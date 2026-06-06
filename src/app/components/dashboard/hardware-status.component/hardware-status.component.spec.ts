import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HardwareStatusComponent } from './hardware-status.component';

describe('HardwareStatusComponent', () => {
  let component: HardwareStatusComponent;
  let fixture: ComponentFixture<HardwareStatusComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HardwareStatusComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(HardwareStatusComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
