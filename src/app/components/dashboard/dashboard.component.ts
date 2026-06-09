import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MapViewComponent } from './map-view.component/map-view.component';
import { ATCCComponent } from '../atcc/atcc.component';
import { DashboardDataService, MenuItem } from '../../services/dashboard-data.service';
import { AtccSection } from '../../services/atcc-data.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
  imports: [
    CommonModule,
    MapViewComponent,
    ATCCComponent,
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

  activeBody = 'dashboard';
  activeAtccSection: AtccSection = 'visualization';
  expandedMenu: string | null = null;
  activeTopTab = 'Dashboard';

  selectMenu(item: MenuItem): void {
    if (!item.children?.length) {
      if (item.route !== undefined) {
        this.activeBody = item.route;
      }
      this.activeTopTab = item.label;
      this.expandedMenu = null;
      return;
    }

    this.expandedMenu = this.expandedMenu === item.label ? null : item.label;
    this.activeTopTab = item.label;
  }

  selectLeaf(item: MenuItem): void {
    if (!item.route) {
      return;
    }

    this.activeBody = item.route;
    this.activeTopTab = 'ATCC';

    if (item.route.startsWith('atcc/')) {
      const section = item.route.split('/')[1] as AtccSection | undefined;
      this.activeAtccSection = section ?? 'visualization';
      this.expandedMenu = 'ATCC';
    }
  }

  hasActiveLeaf(item: MenuItem): boolean {
    if (item.route && item.route === this.activeBody) {
      return true;
    }

    return (item.children ?? []).some(child => this.hasActiveLeaf(child));
  }

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

  sendPtzCommand(cameraId: string, command: 'left' | 'right' | 'up' | 'down' | 'zoomIn' | 'zoomOut'): void {
    void fetch(`/api/ptz/${encodeURIComponent(cameraId)}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ command }),
    }).catch(() => undefined);
  }

}
