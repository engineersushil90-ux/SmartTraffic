import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TrafficStatsComponent } from './traffic-stats.component';

describe('TrafficStatsComponent', () => {
  let component: TrafficStatsComponent;
  let fixture: ComponentFixture<TrafficStatsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TrafficStatsComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(TrafficStatsComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
