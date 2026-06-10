import { AfterViewInit, ChangeDetectorRef, Component, ElementRef, OnDestroy, ViewChild, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import * as L from 'leaflet';

interface CameraSite {
  label: string;
  cameraId: string;
  status: 'Connected' | 'Disconnected';
  stream: string;
  lastSeen: string;
  position: L.LatLngExpression;
  screenX: number;
  screenY: number;
}

@Component({
  selector: 'app-map-view',
  templateUrl: './map-view.component.html',
  styleUrls: ['./map-view.component.scss'],
  standalone: true,
  imports: [CommonModule],
})
export class MapViewComponent implements AfterViewInit, OnDestroy {
  @ViewChild('mapContainer') private mapContainer?: ElementRef<HTMLDivElement>;

  private readonly cdr = inject(ChangeDetectorRef);

  public isBrowser = true;
  activeCamera: CameraSite | null = null;
  readonly cameras: CameraSite[] = [
    {
      label: 'Rohini',
      cameraId: 'CAM-ROH-01',
      status: 'Connected',
      stream: 'FLV live feed',
      lastSeen: 'Live',
      position: [28.7357, 77.0828],
      screenX: 0,
      screenY: 0,
    },
    {
      label: 'Wazirpur',
      cameraId: 'CAM-WAZ-01',
      status: 'Disconnected',
      stream: 'MJPEG feed',
      lastSeen: 'Stream unavailable',
      position: [28.6995, 77.1657],
      screenX: 0,
      screenY: 0,
    },
    {
      label: 'New Delhi',
      cameraId: 'CAM-ND-01',
      status: 'Connected',
      stream: 'Traffic camera site',
      lastSeen: 'Live',
      position: [28.6139, 77.2090],
      screenX: 0,
      screenY: 0,
    },
  ];

  private map?: L.Map;
  private resizeObserver?: ResizeObserver;

  ngAfterViewInit(): void {
    setTimeout(() => this.initMap(), 0);
  }

  ngOnDestroy(): void {
    this.resizeObserver?.disconnect();
    this.map?.remove();
  }

  private initMap(): void {
    const container = this.mapContainer?.nativeElement;

    if (!container || this.map) {
      return;
    }

    this.map = L.map(container, {
      center: [28.6139, 77.2090],
      zoom: 11,
      zoomControl: true,
    });

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors',
      maxZoom: 19,
    }).addTo(this.map);

    this.map.on('click', event => {
      this.activeCamera = null;
      L.popup()
        .setLatLng(event.latlng)
        .setContent(`Map location<br>${event.latlng.lat.toFixed(5)}, ${event.latlng.lng.toFixed(5)}`)
        .openOn(this.map!);
    });

    this.map.on('move zoom resize', () => this.updateCameraPositions());

    this.resizeObserver = new ResizeObserver(() => this.map?.invalidateSize());
    this.resizeObserver.observe(this.map.getContainer());
    [150, 500, 1000].forEach(delay => {
      setTimeout(() => {
        this.map?.invalidateSize();
        this.updateCameraPositions();
      }, delay);
    });
    this.updateCameraPositions();
  }

  showCameraDetails(camera: CameraSite, event: MouseEvent): void {
    event.stopPropagation();
    this.activeCamera = camera;
  }

  private updateCameraPositions(): void {
    if (!this.map) {
      return;
    }

    this.cameras.forEach(camera => {
      const point = this.map!.latLngToContainerPoint(L.latLng(camera.position));
      camera.screenX = point.x;
      camera.screenY = point.y;
    });
    this.cdr.detectChanges();
  }
}
