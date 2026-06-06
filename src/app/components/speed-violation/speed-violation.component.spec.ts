import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SpeedViolation } from './speed-violation.component';

describe('SpeedViolation', () => {
  let component: SpeedViolation;
  let fixture: ComponentFixture<SpeedViolation>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SpeedViolation],
    }).compileComponents();

    fixture = TestBed.createComponent(SpeedViolation);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
