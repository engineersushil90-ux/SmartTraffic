import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CongestionStatsComponent } from './congestion-stats.component';

describe('CongestionStatsComponent', () => {
  let component: CongestionStatsComponent;
  let fixture: ComponentFixture<CongestionStatsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CongestionStatsComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(CongestionStatsComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
