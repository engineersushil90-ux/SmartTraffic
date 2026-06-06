import { Component } from "@angular/core";

// atcc.component.ts
@Component({
  selector: 'app-atcc',
  templateUrl: './atcc.component.html',
  styleUrls: ['./atcc.component.scss'],
  standalone: false
})
export class ATCCComponent {
  intersections = [
    { id: 'Int-01', delay: 24, occupancy: 82, pedestrians: 145 },
    { id: 'Int-02', delay: 39, occupancy: 76, pedestrians: 210 }
  ];
}
