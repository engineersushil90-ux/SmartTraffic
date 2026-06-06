import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IncidentSummaryComponent } from './incident-summary.component';

describe('IncidentSummaryComponent', () => {
  let component: IncidentSummaryComponent;
  let fixture: ComponentFixture<IncidentSummaryComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [IncidentSummaryComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(IncidentSummaryComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
