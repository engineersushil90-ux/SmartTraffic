import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MapViewComponent } from './map-view.component/map-view.component';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
  imports: [
    CommonModule,
    MapViewComponent,
  ],
})
export class DashboardComponent {
  isLightTheme = false;
  isSidebarHidden = false;
  readonly notificationCount = 8;
  readonly userName = 'admin';

  toggleSidebar(): void {
    this.isSidebarHidden = !this.isSidebarHidden;
  }

  toggleTheme(): void {
    this.isLightTheme = !this.isLightTheme;
  }

  toggleFullscreen(): void {
    const documentElement = document.documentElement;

    if (!document.fullscreenElement) {
      documentElement.requestFullscreen?.();
      return;
    }

    document.exitFullscreen?.();
  }
}
