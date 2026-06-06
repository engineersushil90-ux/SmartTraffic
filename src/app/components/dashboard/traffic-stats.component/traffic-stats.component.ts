import { Component } from '@angular/core';
import { MatCard } from "@angular/material/card";

@Component({
  selector: 'app-traffic-stats',
  templateUrl: './traffic-stats.component.html',
  styleUrl: './traffic-stats.component.scss',
  imports: [MatCard],
})
export class TrafficStatsComponent {
  trafficStats = {
    avgSpeed: 29,
    volume: 36200,
    alerts: 5,
    congestionIndex: 2.6
  };
}
