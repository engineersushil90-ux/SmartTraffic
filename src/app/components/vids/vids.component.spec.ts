import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Vids } from './vids.component';

describe('Vids', () => {
  let component: Vids;
  let fixture: ComponentFixture<Vids>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Vids],
    }).compileComponents();

    fixture = TestBed.createComponent(Vids);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
