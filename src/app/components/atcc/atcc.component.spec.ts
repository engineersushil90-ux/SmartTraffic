import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Atcc } from './atcc.component';

describe('Atcc', () => {
  let component: Atcc;
  let fixture: ComponentFixture<Atcc>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Atcc],
    }).compileComponents();

    fixture = TestBed.createComponent(Atcc);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
