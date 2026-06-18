import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VIDSComponent } from './vids.component';

describe('VIDSComponent', () => {
  let component: VIDSComponent;
  let fixture: ComponentFixture<VIDSComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [VIDSComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(VIDSComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
