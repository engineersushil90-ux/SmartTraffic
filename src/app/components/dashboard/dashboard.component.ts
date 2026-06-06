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
  readonly sideMenuItems = [
    ['H', 'Dashboard'], ['A', 'ATCC'], ['V', 'VIDS'], ['M', 'VMS'], ['W', 'MET'], ['S', 'VSDS'],
    ['I', 'Incidents'], ['!', 'Alerts'], ['T', 'Traffic Explorer'], ['G', 'Trend Analytics'],
    ['C', 'Device Management'], ['S', 'System Health'], ['U', 'User Management'], ['L', 'Audit Logs'],
  ];
  readonly topTabs = ['Dashboard', 'ATCC', 'VIDS', 'VMS', 'MET', 'VSDS', 'Reports', 'Analytics', 'Settings'];
  readonly kpis = [
    { label: 'Average Speed', value: '42', unit: 'km/h', trend: '5% vs yesterday', state: 'up', icon: 'speed-icon' },
    { label: 'Traffic Volume', value: '36,200', unit: 'veh/hr', trend: '8.4% vs yesterday', state: 'up', icon: 'vehicle-icon' },
    { label: 'Active Incidents', value: '12', unit: '', trend: '2 new', state: 'down', icon: 'warning-icon' },
    { label: 'Travel Time Index', value: '2.6', unit: '', trend: 'Moderate', state: 'moderate', icon: 'sparkline' },
    { label: 'Congestion Level', value: 'Moderate', unit: '', trend: '', state: '', meter: true },
    { label: 'Weather', value: '28', unit: 'C', trend: 'Partly Cloudy', state: '', icon: 'weather-icon' },
  ];
  readonly mapLegend = [
    ['smooth', 'Smooth'], ['moderate-dot', 'Moderate'], ['congested', 'Congested'], ['incident', 'Incident'], ['camera', 'Camera'],
  ];
  readonly videoFeeds = ['NH 44 - Rohini', 'Ring Road - Wazirpur', 'NH 48 - Rajouri Garden', 'NH 24 - DND Flyway'];
  readonly alerts = [
    { color: 'red', icon: '!', text: 'Accident on NH 48', time: '10:40 AM' },
    { color: 'orange', icon: '!', text: 'Heavy Traffic on Ring Road', time: '10:38 AM' },
    { color: 'yellow', icon: '!', text: 'Lane Closure on NH 24', time: '10:35 AM' },
    { color: 'blue', icon: '', text: 'Poor Visibility - Fog', time: '10:30 AM' },
    { color: 'purple', icon: 'A', text: 'Signal Malfunction - Sector 62', time: '10:28 AM' },
    { color: 'orange', icon: '!', text: 'Queue build-up at Toll Plaza', time: '10:24 AM' },
    { color: 'blue', icon: '', text: 'Camera offline at Wazirpur', time: '10:21 AM' },
  ];
  readonly vmsMessages = [
    { message: 'DRIVE SAFE', line: 'ARRIVE SAFE', place: 'IIT Flyover', time: '10:42 AM', status: 'Active' },
    { message: 'HEAVY TRAFFIC', line: 'AHEAD', place: 'NH 44 Rohini', time: '10:41 AM', status: 'Active' },
    { message: 'FOG AHEAD', line: 'DRIVE SLOW', place: 'DND Flyway', time: '10:40 AM', status: 'Active' },
    { message: 'SLOW DOWN', line: 'WORK ZONE', place: 'Ring Road', time: '10:36 AM', status: 'Active' },
  ];
  readonly weatherStats = [
    ['Humidity', '54%'], ['Wind', '12 km/h NE'], ['Rainfall', '0.0 mm'], ['Visibility', '6.5 km'],
  ];
  readonly vsdsMetrics = [
    { label: 'Vehicle Count', value: '36,200', trend: '8.4% vs yesterday', state: 'up' },
    { label: 'Average Speed', value: '95.7 km/h', trend: '2.1% vs yesterday', state: 'down' },
    { label: 'Maximum Speed', value: '213 km/h', trend: 'Highest today', state: 'moderate' },
    { label: 'Speed Violations', value: '36,097', trend: '12% vs yesterday', state: 'up' },
  ];
  readonly incidents = [
    { color: 'red', name: 'Accident', count: 5, percent: 42 },
    { color: 'orange', name: 'Breakdown', count: 3, percent: 25 },
    { color: 'yellow', name: 'Construction', count: 2, percent: 17 },
    { color: 'gray', name: 'Other', count: 2, percent: 16 },
    { color: 'blue', name: 'Camera Issue', count: 1, percent: 8 },
  ];
  readonly deviceSummary = [
    ['Total Device', '85'], ['Connected', '66'], ['Disconnected', '19'],
  ];
  readonly deviceStatus = [
    { name: 'ATCC', online: 0, offline: 6 }, { name: 'ECB-GSM', online: 28, offline: 0 },
    { name: 'ECB-SYS', online: 0, offline: 1 }, { name: 'Main', online: 1, offline: 0 },
    { name: 'MET', online: 1, offline: 0 }, { name: 'PTZ', online: 16, offline: 6 },
    { name: 'VIDS', online: 10, offline: 0 }, { name: 'VMS', online: 2, offline: 6 },
    { name: 'VSDS', online: 6, offline: 0 },
  ];
  readonly activities = [
    ['10:43 AM', 'Incident created on NH 48'], ['10:41 AM', 'VMS message updated on NH 44'],
    ['10:40 AM', 'Camera offline at Sector 62'], ['10:38 AM', 'Weather alert: Fog in Delhi'],
    ['10:35 AM', 'Traffic data synced from VSDS-101'], ['10:31 AM', 'PTZ camera recovered at Rohini'],
    ['10:28 AM', 'Alert acknowledged by admin'],
  ];

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
