import { Injectable } from '@angular/core';

export type TrendState = 'up' | 'down' | 'moderate' | '';

export interface MenuItem {
  icon: string;
  label: string;
}

export interface KpiCard {
  label: string;
  value: string;
  unit: string;
  trend: string;
  state: TrendState;
  icon?: string;
  meter?: boolean;
}

export interface LegendItem {
  className: string;
  label: string;
}

export interface AlertItem {
  color: string;
  icon: string;
  text: string;
  time: string;
}

export interface VmsMessage {
  message: string;
  line: string;
  place: string;
  time: string;
  status: string;
}

export interface NameValue {
  label: string;
  value: string;
}

export interface VsdsMetric {
  label: string;
  value: string;
  trend: string;
  state: TrendState;
}

export interface IncidentMetric {
  color: string;
  name: string;
  count: number;
  percent: number;
}

export interface DeviceHealthMetric {
  label: string;
  value: number;
}

export interface DeviceStatusMetric {
  name: string;
  online: number;
  offline: number;
}

export interface ActivityItem {
  time: string;
  text: string;
}

export interface DashboardData {
  header: {
    notificationCount: number;
    userName: string;
    clock: string;
  };
  navigation: {
    sideMenuItems: MenuItem[];
    topTabs: string[];
  };
  kpis: KpiCard[];
  mapLegend: LegendItem[];
  videoFeeds: string[];
  alerts: AlertItem[];
  vmsMessages: VmsMessage[];
  weatherStats: NameValue[];
  vsdsMetrics: VsdsMetric[];
  incidents: IncidentMetric[];
  deviceSummary: NameValue[];
  deviceHealthMetrics: DeviceHealthMetric[];
  deviceStatus: DeviceStatusMetric[];
  activities: ActivityItem[];
}

@Injectable({
  providedIn: 'root',
})
export class DashboardDataService {
  getDashboardData(): DashboardData {
    return {
      header: {
        notificationCount: 8,
        userName: 'admin',
        clock: 'Tue, Apr 24 | 10:45 AM',
      },
      navigation: {
        sideMenuItems: [
          { icon: 'H', label: 'Dashboard' },
          { icon: 'A', label: 'ATCC' },
          { icon: 'V', label: 'VIDS' },
          { icon: 'M', label: 'VMS' },
          { icon: 'W', label: 'MET' },
          { icon: 'S', label: 'VSDS' },
          { icon: 'I', label: 'Incidents' },
          { icon: '!', label: 'Alerts' },
          { icon: 'T', label: 'Traffic Explorer' },
          { icon: 'G', label: 'Trend Analytics' },
          { icon: 'C', label: 'Device Management' },
          { icon: 'S', label: 'System Health' },
          { icon: 'U', label: 'User Management' },
          { icon: 'L', label: 'Audit Logs' },
        ],
        topTabs: ['Dashboard', 'ATCC', 'VIDS', 'VMS', 'MET', 'VSDS', 'Reports', 'Analytics', 'Settings'],
      },
      kpis: [
        { label: 'Average Speed', value: '42', unit: 'km/h', trend: '5% vs yesterday', state: 'up', icon: 'speed-icon' },
        { label: 'Traffic Volume', value: '36,200', unit: 'veh/hr', trend: '8.4% vs yesterday', state: 'up', icon: 'vehicle-icon' },
        { label: 'Active Incidents', value: '12', unit: '', trend: '2 new', state: 'down', icon: 'warning-icon' },
        { label: 'Travel Time Index', value: '2.6', unit: '', trend: 'Moderate', state: 'moderate', icon: 'sparkline' },
        { label: 'Congestion Level', value: 'Moderate', unit: '', trend: '', state: '', meter: true },
        { label: 'Weather', value: '28', unit: 'C', trend: 'Partly Cloudy', state: '', icon: 'weather-icon' },
      ],
      mapLegend: [
        { className: 'smooth', label: 'Smooth' },
        { className: 'moderate-dot', label: 'Moderate' },
        { className: 'congested', label: 'Congested' },
        { className: 'incident', label: 'Incident' },
        { className: 'camera', label: 'Camera' },
      ],
      videoFeeds: ['NH 44 - Rohini', 'Ring Road - Wazirpur', 'NH 48 - Rajouri Garden', 'NH 24 - DND Flyway'],
      alerts: [
        { color: 'red', icon: '!', text: 'Accident on NH 48', time: '10:40 AM' },
        { color: 'orange', icon: '!', text: 'Heavy Traffic on Ring Road', time: '10:38 AM' },
        { color: 'yellow', icon: '!', text: 'Lane Closure on NH 24', time: '10:35 AM' },
        { color: 'blue', icon: '', text: 'Poor Visibility - Fog', time: '10:30 AM' },
        { color: 'purple', icon: 'A', text: 'Signal Malfunction - Sector 62', time: '10:28 AM' },
        { color: 'orange', icon: '!', text: 'Queue build-up at Toll Plaza', time: '10:24 AM' },
        { color: 'blue', icon: '', text: 'Camera offline at Wazirpur', time: '10:21 AM' },
      ],
      vmsMessages: [
        { message: 'DRIVE SAFE', line: 'ARRIVE SAFE', place: 'IIT Flyover', time: '10:42 AM', status: 'Active' },
        { message: 'HEAVY TRAFFIC', line: 'AHEAD', place: 'NH 44 Rohini', time: '10:41 AM', status: 'Active' },
        { message: 'FOG AHEAD', line: 'DRIVE SLOW', place: 'DND Flyway', time: '10:40 AM', status: 'Active' },
        { message: 'SLOW DOWN', line: 'WORK ZONE', place: 'Ring Road', time: '10:36 AM', status: 'Active' },
      ],
      weatherStats: [
        { label: 'Humidity', value: '54%' },
        { label: 'Wind', value: '12 km/h NE' },
        { label: 'Rainfall', value: '0.0 mm' },
        { label: 'Visibility', value: '6.5 km' },
      ],
      vsdsMetrics: [
        { label: 'Vehicle Count', value: '36,200', trend: '8.4% vs yesterday', state: 'up' },
        { label: 'Average Speed', value: '95.7 km/h', trend: '2.1% vs yesterday', state: 'down' },
        { label: 'Maximum Speed', value: '213 km/h', trend: 'Highest today', state: 'moderate' },
        { label: 'Speed Violations', value: '36,097', trend: '12% vs yesterday', state: 'up' },
      ],
      incidents: [
        { color: 'red', name: 'Accident', count: 5, percent: 42 },
        { color: 'orange', name: 'Breakdown', count: 3, percent: 25 },
        { color: 'yellow', name: 'Construction', count: 2, percent: 17 },
        { color: 'gray', name: 'Other', count: 2, percent: 16 },
        { color: 'blue', name: 'Camera Issue', count: 1, percent: 8 },
      ],
      deviceSummary: [
        { label: 'Total Device', value: '85' },
        { label: 'Connected', value: '66' },
        { label: 'Disconnected', value: '19' },
      ],
      deviceHealthMetrics: [
        { label: 'Device Health', value: 92 },
        { label: 'Road Occupancy', value: 24 },
        { label: 'Compliance Rate', value: 88 },
        { label: 'Camera Availability', value: 97 },
        { label: 'Network Connectivity', value: 99 },
      ],
      deviceStatus: [
        { name: 'ATCC', online: 0, offline: 6 },
        { name: 'ECB-GSM', online: 28, offline: 0 },
        { name: 'ECB-SYS', online: 0, offline: 1 },
        { name: 'Main', online: 1, offline: 0 },
        { name: 'MET', online: 1, offline: 0 },
        { name: 'PTZ', online: 16, offline: 6 },
        { name: 'VIDS', online: 10, offline: 0 },
        { name: 'VMS', online: 2, offline: 6 },
        { name: 'VSDS', online: 6, offline: 0 },
      ],
      activities: [
        { time: '10:43 AM', text: 'Incident created on NH 48' },
        { time: '10:41 AM', text: 'VMS message updated on NH 44' },
        { time: '10:40 AM', text: 'Camera offline at Sector 62' },
        { time: '10:38 AM', text: 'Weather alert: Fog in Delhi' },
        { time: '10:35 AM', text: 'Traffic data synced from VSDS-101' },
        { time: '10:31 AM', text: 'PTZ camera recovered at Rohini' },
        { time: '10:28 AM', text: 'Alert acknowledged by admin' },
      ],
    };
  }
}
