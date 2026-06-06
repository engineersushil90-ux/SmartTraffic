import { Component, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';   // ✅ import CommonModule
import * as L from 'leaflet';

@Component({
  selector: 'app-map-view',
  templateUrl: './map-view.component.html',
  styleUrls: ['./map-view.component.scss'],
  standalone: true,
  imports: [CommonModule]   // ✅ add CommonModule here
})
export class MapViewComponent implements AfterViewInit {
  public isBrowser = true;
  private map!: L.Map;

  ngAfterViewInit(): void {
     setTimeout(() => this.initMap(), 0);
  }

  private initMap(): void {
    this.map = L.map('map', {
      center: [28.4595, 77.0266], // Gurgaon
      zoom: 12
    });

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.map);

    const marker = L.marker([28.467, 77.030]).addTo(this.map);
    marker.bindPopup('<b>Traffic Camera</b><br>Sector 8, Gurgaon').openPopup();
  }
}
