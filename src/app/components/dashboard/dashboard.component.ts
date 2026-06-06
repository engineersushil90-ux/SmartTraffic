import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MapViewComponent } from './map-view.component/map-view.component';
import { DashboardDataService } from '../../services/dashboard-data.service';

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
  private readonly dashboardDataService = inject(DashboardDataService);
  private readonly dashboardData = this.dashboardDataService.getDashboardData();

  isLightTheme = false;
  isSidebarHidden = false;

  readonly notificationCount = this.dashboardData.header.notificationCount;
  readonly userName = this.dashboardData.header.userName;
  readonly clock = this.dashboardData.header.clock;
  readonly sideMenuItems = this.dashboardData.navigation.sideMenuItems;
  readonly topTabs = this.dashboardData.navigation.topTabs;
  readonly kpis = this.dashboardData.kpis;
  readonly mapLegend = this.dashboardData.mapLegend;
  readonly videoFeeds = this.dashboardData.videoFeeds;
  readonly alerts = this.dashboardData.alerts;
  readonly vmsMessages = this.dashboardData.vmsMessages;
  readonly weatherStats = this.dashboardData.weatherStats;
  readonly vsdsMetrics = this.dashboardData.vsdsMetrics;
  readonly incidents = this.dashboardData.incidents;
  readonly deviceSummary = this.dashboardData.deviceSummary;
  readonly deviceHealthMetrics = this.dashboardData.deviceHealthMetrics;
  readonly deviceStatus = this.dashboardData.deviceStatus;
  readonly activities = this.dashboardData.activities;

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
