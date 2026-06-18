import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { Component, Input, OnChanges, OnInit, SimpleChanges, inject } from '@angular/core';

export type VidsSection = 'visualization' | 'table' | 'reports';

interface VidsSummaryCard {
  label: string;
  value: string;
  meta: string;
  icon: string;
  tone: string;
}

interface CameraFeed {
  id: string;
  name: string;
  scene: string;
  status: 'Live' | 'Review';
  detections: Array<{ label: string; tone: string; left: number; top: number; width: number; height: number }>;
}

interface IncidentRow {
  id: number;
  time: string;
  type: string;
  camera: string;
  location: string;
  severity: 'High' | 'Medium' | 'Low';
  status: 'Active' | 'Resolved';
  confidence: string;
  tone: string;
}

interface ChartMetric {
  label: string;
  value: number;
  color: string;
}

interface ReportParameter {
  label: string;
  type: 'date' | 'time' | 'select';
  value: string;
  options?: string[];
}

interface ReportCard {
  title: string;
  description: string;
  parameters: ReportParameter[];
  exports: string[];
}

@Component({
  selector: 'app-vids',
  imports: [CommonModule],
  templateUrl: './vids.component.html',
  styleUrl: './vids.component.scss',
})
export class VIDSComponent implements OnInit, OnChanges {
  @Input() embedded = false;
  @Input() section: VidsSection | null = null;

  private readonly route = inject(ActivatedRoute);

  readonly updatedAt = '10:42:10 AM';
  readonly summaryCards: VidsSummaryCard[] = [
    { label: 'Total Cameras', value: '128', meta: 'Online 118', icon: 'CAM', tone: 'blue' },
    { label: 'AI Events Today', value: '346', meta: '15.2% vs yesterday', icon: 'AI', tone: 'cyan' },
    { label: 'Active Incidents', value: '12', meta: 'High 5  Medium 7', icon: '!', tone: 'red' },
    { label: 'Detection Accuracy', value: '98.7%', meta: '2.3% vs yesterday', icon: 'OK', tone: 'indigo' },
    { label: 'Avg Response Time', value: '2m 14s', meta: '18s vs yesterday', icon: 'CLK', tone: 'green' },
    { label: 'False Positives', value: '3', meta: '25% vs yesterday', icon: 'FP', tone: 'gold' },
  ];

  readonly cameraFeeds: CameraFeed[] = [
    {
      id: 'CAM 01',
      name: 'MG Road',
      scene: 'mg-road',
      status: 'Live',
      detections: [
        { label: 'CAR', tone: 'blue', left: 12, top: 34, width: 12, height: 14 },
        { label: 'PERSON', tone: 'green', left: 48, top: 35, width: 10, height: 16 },
        { label: 'TRUCK', tone: 'blue', left: 70, top: 24, width: 10, height: 14 },
      ],
    },
    {
      id: 'CAM 02',
      name: 'Ring Road',
      scene: 'ring-road',
      status: 'Live',
      detections: [
        { label: 'WRONG WAY', tone: 'red', left: 43, top: 28, width: 14, height: 17 },
        { label: 'CAR', tone: 'blue', left: 22, top: 47, width: 16, height: 17 },
      ],
    },
    {
      id: 'CAM 03',
      name: 'City Center',
      scene: 'city-center',
      status: 'Live',
      detections: [
        { label: 'PERSON', tone: 'green', left: 8, top: 22, width: 8, height: 20 },
        { label: 'CAR', tone: 'blue', left: 37, top: 42, width: 14, height: 17 },
        { label: 'AUTO', tone: 'purple', left: 73, top: 39, width: 12, height: 15 },
      ],
    },
    {
      id: 'CAM 04',
      name: 'Highway Entry',
      scene: 'highway-entry',
      status: 'Live',
      detections: [
        { label: 'CAR', tone: 'blue', left: 20, top: 43, width: 14, height: 17 },
        { label: 'STOPPED VEHICLE', tone: 'red', left: 56, top: 39, width: 18, height: 20 },
      ],
    },
  ];

  readonly incidents: IncidentRow[] = [
    { id: 1, time: '10:42:10 AM', type: 'Wrong Way Vehicle', camera: 'CAM 02', location: 'Ring Road', severity: 'High', status: 'Active', confidence: '96.4%', tone: 'red' },
    { id: 2, time: '10:39:48 AM', type: 'Pedestrian Crossing', camera: 'CAM 03', location: 'City Center', severity: 'Medium', status: 'Active', confidence: '92.1%', tone: 'gold' },
    { id: 3, time: '10:36:21 AM', type: 'Stopped Vehicle', camera: 'CAM 04', location: 'Highway Entry', severity: 'Medium', status: 'Active', confidence: '95.3%', tone: 'gold' },
    { id: 4, time: '10:32:05 AM', type: 'Lane Violation', camera: 'CAM 01', location: 'MG Road', severity: 'Low', status: 'Active', confidence: '90.7%', tone: 'blue' },
    { id: 5, time: '10:28:11 AM', type: 'Road Obstruction', camera: 'CAM 05', location: 'Bypass Road', severity: 'Medium', status: 'Resolved', confidence: '93.2%', tone: 'blue' },
    { id: 6, time: '10:21:34 AM', type: 'Accident Detected', camera: 'CAM 06', location: 'Outer Ring Rd', severity: 'High', status: 'Active', confidence: '97.6%', tone: 'red' },
    { id: 7, time: '10:18:47 AM', type: 'Stopped Vehicle', camera: 'CAM 02', location: 'Ring Road', severity: 'Medium', status: 'Resolved', confidence: '94.0%', tone: 'gold' },
    { id: 8, time: '10:14:22 AM', type: 'Pedestrian Detection', camera: 'CAM 03', location: 'City Center', severity: 'Medium', status: 'Resolved', confidence: '91.8%', tone: 'gold' },
    { id: 9, time: '10:09:15 AM', type: 'Lane Violation', camera: 'CAM 01', location: 'MG Road', severity: 'Low', status: 'Resolved', confidence: '89.2%', tone: 'blue' },
    { id: 10, time: '10:05:31 AM', type: 'Wrong Way Vehicle', camera: 'CAM 04', location: 'Highway Entry', severity: 'High', status: 'Resolved', confidence: '96.8%', tone: 'red' },
  ];

  readonly incidentTypes: ChartMetric[] = [
    { label: 'Wrong Way', value: 28, color: '#e25844' },
    { label: 'Stopped Vehicle', value: 22, color: '#f0b642' },
    { label: 'Pedestrian', value: 18, color: '#68c46c' },
    { label: 'Lane Violation', value: 15, color: '#4f85e8' },
    { label: 'Accident', value: 10, color: '#765bd5' },
    { label: 'Others', value: 7, color: '#9aa4b5' },
  ];

  readonly cameraStatus = [
    { label: 'Online', value: '118 (92%)', color: '#67c56a' },
    { label: 'Offline', value: '6 (5%)', color: '#e25a45' },
    { label: 'Maintenance', value: '4 (3%)', color: '#f1b847' },
  ];

  readonly reports: ReportCard[] = [
    {
      title: 'Incident Summary Report',
      description: 'Aggregated VIDS incident counts, severity distribution, and response performance.',
      parameters: [
        { label: 'Start Date', type: 'date', value: '2026-06-01' },
        { label: 'End Date', type: 'date', value: '2026-06-18' },
        { label: 'Start Time', type: 'time', value: '00:00' },
        { label: 'End Time', type: 'time', value: '23:59:59' },
        { label: 'Location', type: 'select', value: 'All Locations', options: ['All Locations', 'MG Road', 'Ring Road', 'City Center', 'Highway Entry'] },
        { label: 'Severity', type: 'select', value: 'All Severities', options: ['All Severities', 'High', 'Medium', 'Low'] },
        { label: 'Status', type: 'select', value: 'All Statuses', options: ['All Statuses', 'Active', 'Resolved'] },
        { label: 'Camera ID', type: 'select', value: 'All Cameras', options: ['All Cameras', 'CAM 01', 'CAM 02', 'CAM 03', 'CAM 04'] },
      ],
      exports: ['PDF', 'Excel', 'CSV'],
    },
    {
      title: 'Detailed Incident Report',
      description: 'Record-level VIDS incident log with camera, timestamp, severity, status, and AI confidence.',
      parameters: [
        { label: 'Start Date', type: 'date', value: '2026-06-01' },
        { label: 'End Date', type: 'date', value: '2026-06-18' },
        { label: 'Incident Type', type: 'select', value: 'All Types', options: ['All Types', 'Wrong Way Vehicle', 'Stopped Vehicle', 'Lane Violation', 'Accident Detected'] },
        { label: 'AI Confidence', type: 'select', value: '90% and above', options: ['All', '90% and above', '95% and above'] },
        { label: 'Location', type: 'select', value: 'All Locations', options: ['All Locations', 'MG Road', 'Ring Road', 'City Center', 'Highway Entry'] },
        { label: 'Camera ID', type: 'select', value: 'All Cameras', options: ['All Cameras', 'CAM 01', 'CAM 02', 'CAM 03', 'CAM 04'] },
      ],
      exports: ['PDF', 'Excel', 'CSV'],
    },
    {
      title: 'Report with Evidence',
      description: 'VIDS incidents with visual evidence snapshots and operator review status.',
      parameters: [
        { label: 'Start Date', type: 'date', value: '2026-06-01' },
        { label: 'End Date', type: 'date', value: '2026-06-18' },
        { label: 'Evidence Type', type: 'select', value: 'Snapshots and Clips', options: ['Snapshots and Clips', 'Snapshots Only', 'Clips Only'] },
        { label: 'Status', type: 'select', value: 'All Statuses', options: ['All Statuses', 'Active', 'Resolved'] },
        { label: 'Reviewer', type: 'select', value: 'All Reviewers', options: ['All Reviewers', 'Control Room', 'Supervisor'] },
        { label: 'Camera ID', type: 'select', value: 'All Cameras', options: ['All Cameras', 'CAM 01', 'CAM 02', 'CAM 03', 'CAM 04'] },
      ],
      exports: ['PDF', 'Excel', 'CSV'],
    },
  ];

  readonly subMenu: Array<{ id: VidsSection; label: string; icon: string }> = [
    { id: 'visualization', label: 'Visualization', icon: 'view' },
    { id: 'table', label: 'Table', icon: 'grid' },
    { id: 'reports', label: 'Reports', icon: 'file' },
  ];

  readonly sectionSubMenu: Record<VidsSection, Array<{ id: string; label: string }>> = {
    visualization: [
      { id: 'live-view', label: 'Live View' },
      { id: 'analytics', label: 'Analytics' },
      { id: 'cameras', label: 'Cameras' },
    ],
    table: [
      { id: 'incident-log', label: 'Incident Log' },
      { id: 'active', label: 'Active' },
      { id: 'resolved', label: 'Resolved' },
    ],
    reports: [
      { id: 'summary', label: 'Summary' },
      { id: 'detailed', label: 'Detailed' },
      { id: 'evidence', label: 'Evidence' },
    ],
  };

  activeSection: VidsSection = 'visualization';
  activeSubSection = this.sectionSubMenu[this.activeSection][0].id;
  selectedReportIndex = 0;
  selectedCameraIds = this.cameraFeeds.slice(0, 4).map(camera => camera.id);

  get selectedCameraFeeds(): CameraFeed[] {
    const selectedIds = new Set(this.selectedCameraIds);
    return this.cameraFeeds.filter(camera => selectedIds.has(camera.id));
  }

  ngOnInit(): void {
    if (this.section && this.subMenu.some(item => item.id === this.section)) {
      this.setSection(this.section);
    }

    if (!this.embedded) {
      this.route.paramMap.subscribe(params => {
        const section = params.get('section') as VidsSection | null;
        this.setSection(section && this.subMenu.some(item => item.id === section) ? section : 'visualization');
      });
    }
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['section'] && this.embedded && changes['section'].currentValue) {
      const section = changes['section'].currentValue as VidsSection;
      if (this.subMenu.some(item => item.id === section)) {
        this.setSection(section);
      }
    }
  }

  get selectedReport(): ReportCard {
    return this.reports[this.selectedReportIndex] ?? this.reports[0];
  }

  setSection(section: VidsSection): void {
    this.activeSection = section;
    this.activeSubSection = this.sectionSubMenu[section][0].id;
  }

  setSubSection(subSection: string): void {
    this.activeSubSection = subSection;
  }

  isCameraSelected(cameraId: string): boolean {
    return this.selectedCameraIds.includes(cameraId);
  }

  setSelectedCameras(event: Event): void {
    const select = event.target as HTMLSelectElement;
    const selectedIds = Array.from(select.selectedOptions).map(option => option.value).slice(0, 4);
    this.selectedCameraIds = selectedIds.length ? selectedIds : [this.cameraFeeds[0].id];

    Array.from(select.options).forEach(option => {
      option.selected = this.selectedCameraIds.includes(option.value);
    });
  }

  setReport(index: string): void {
    const nextIndex = Number(index);
    if (Number.isInteger(nextIndex) && this.reports[nextIndex]) {
      this.selectedReportIndex = nextIndex;
    }
  }
}
