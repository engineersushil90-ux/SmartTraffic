import { Injectable } from '@angular/core';

export type AtccSection = 'visualization' | 'table' | 'reports';

export interface AtccSummaryCard {
  label: string;
  value: string;
  meta: string;
  icon: string;
  tone: string;
}

export interface VehicleClassMetric {
  name: string;
  count: number;
  percent: number;
  color: string;
}

export interface DirectionMetric {
  className: string;
  east: number;
  north: number;
  south: number;
  west: number;
}

export interface AtccRecord {
  id: number;
  junction: string;
  timestamp: string;
  vehicleClass: string;
  direction: string;
  cameraId: string;
  lane: number;
  speed: string;
  color: string;
}

export interface ReportCard {
  title: string;
  description: string;
  fields: string[];
  exports: string[];
}

export interface AtccData {
  updatedAt: string;
  summaryCards: AtccSummaryCard[];
  vehicleClasses: VehicleClassMetric[];
  directions: DirectionMetric[];
  records: AtccRecord[];
  reports: ReportCard[];
}

@Injectable({
  providedIn: 'root',
})
export class AtccDataService {
  getAtccData(): AtccData {
    return {
      updatedAt: '9:16:29 pm',
      summaryCards: [
        { label: 'Total Vehicles', value: '15,015', meta: 'All Time', icon: 'car', tone: 'blue' },
        { label: 'Peak Hour', value: '--:--', meta: '0 vehicles', icon: 'clock', tone: 'red' },
        { label: 'Average Hourly', value: '0', meta: 'Vehicles per hour', icon: 'bars', tone: 'indigo' },
        { label: 'Date Range', value: 'All Time', meta: '1/1/2000 to 6/6/2026', icon: 'calendar', tone: 'yellow' },
      ],
      vehicleClasses: [
        { name: 'Sedan', count: 1824, percent: 12.1, color: '#4f5fc9' },
        { name: 'SUV', count: 1806, percent: 12.0, color: '#d8ad48' },
        { name: 'Light Truck', count: 1785, percent: 11.9, color: '#6caf79' },
        { name: 'Motorcycle', count: 1782, percent: 11.9, color: '#db9b45' },
        { name: 'Heavy Truck', count: 1774, percent: 11.8, color: '#5a66cd' },
        { name: 'Bus', count: 1766, percent: 11.8, color: '#5d85d8' },
        { name: 'Auto-rickshaw', count: 1729, percent: 11.5, color: '#df8f8c' },
        { name: 'Bicycle', count: 1606, percent: 10.7, color: '#67a973' },
        { name: 'Truck', count: 751, percent: 5.0, color: '#3e4c9f' },
        { name: 'Van', count: 1, percent: 0.0, color: '#444ca2' },
      ],
      directions: [
        { className: 'Auto', east: 1, north: 1, south: 1, west: 0 },
        { className: 'Auto-rickshaw', east: 380, north: 440, south: 430, west: 479 },
        { className: 'Bike', east: 3, north: 4, south: 2, west: 2 },
        { className: 'Bus', east: 395, north: 455, south: 430, west: 486 },
        { className: 'Bus/Truck', east: 1, north: 1, south: 1, west: 1 },
        { className: 'Car', east: 22, north: 20, south: 23, west: 22 },
        { className: 'Heavy Truck', east: 390, north: 460, south: 430, west: 494 },
        { className: 'Light Truck', east: 1, north: 1, south: 0, west: 1 },
        { className: 'Motorcycle', east: 395, north: 455, south: 430, west: 502 },
        { className: 'SUV', east: 400, north: 465, south: 430, west: 510 },
        { className: 'Truck', east: 390, north: 470, south: 440, west: 485 },
        { className: 'Van', east: 2, north: 2, south: 1, west: 1 },
      ],
      records: [
        { id: 17613, junction: 'ATCC-Test', timestamp: '2026-03-09 18:25:47', vehicleClass: 'Car', direction: 'Toward phil', cameraId: '192.168.2.97', lane: 2, speed: '-', color: '-' },
        { id: 17612, junction: 'ATCC-Test', timestamp: '2026-03-09 18:25:46', vehicleClass: 'Car', direction: 'Toward phil', cameraId: '192.168.2.97', lane: 1, speed: '-', color: '-' },
        { id: 17611, junction: 'ATCC-Test', timestamp: '2026-03-09 18:25:45', vehicleClass: 'Car', direction: 'Toward phil', cameraId: '192.168.2.97', lane: 1, speed: '-', color: '-' },
        { id: 17610, junction: 'ATCC-Test', timestamp: '2026-03-09 18:25:44', vehicleClass: 'Bus/Truck', direction: 'Toward phil', cameraId: '192.168.2.97', lane: 2, speed: '-', color: '-' },
        { id: 17609, junction: 'ATCC-Test', timestamp: '2026-03-09 18:25:43', vehicleClass: 'Car', direction: 'Toward phil', cameraId: '192.168.2.97', lane: 1, speed: '-', color: '-' },
        { id: 17608, junction: 'ATCC-Test', timestamp: '2026-03-09 18:25:43', vehicleClass: 'Car', direction: 'Toward phil', cameraId: '192.168.2.97', lane: 1, speed: '-', color: '-' },
        { id: 17607, junction: 'ATCC-Test', timestamp: '2026-03-09 18:25:41', vehicleClass: 'Bus/Truck', direction: 'Toward phil', cameraId: '192.168.2.97', lane: 2, speed: '-', color: '-' },
      ],
      reports: [
        { title: 'Summary Report', description: 'Overview statistics with total counts, peak hours, and key metrics.', fields: ['Start Date', 'End Date', 'Start Time', 'End Time', 'Junction', 'Vehicle Class'], exports: ['PDF', 'Excel', 'CSV'] },
        { title: 'Detailed Report', description: 'Complete record-by-record listing with timestamps, speeds, and vehicle details.', fields: ['Start Date', 'End Date', 'Start Time', 'End Time', 'Junction', 'Vehicle Class', 'Direction', 'Camera ID'], exports: ['PDF', 'Excel', 'CSV'] },
        { title: 'Report with Images', description: 'Traffic records with vehicle images for visual verification and analysis.', fields: ['Start Date', 'End Date', 'Start Time', 'End Time', 'Junction', 'Vehicle Class', 'Direction', 'Camera ID'], exports: ['PDF', 'Excel', 'CSV'] },
      ],
    };
  }
}
